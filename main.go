package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
)

import (
	"github.com/eyedeekay/eeproxy/socks"
	"github.com/eyedeekay/sam-forwarder/config"
)

type flagOpts []string

func (f *flagOpts) String() string {
	r := ""
	for _, s := range *f {
		r += s + ","
	}
	return strings.TrimSuffix(r, ",")
}

func (f *flagOpts) Set(s string) error {
	*f = append(*f, s)
	return nil
}

func (f *flagOpts) StringSlice() []string {
	var r []string
	for _, s := range *f {
		r = append(r, s)
	}
	return r
}

var (
	startUp = flag.Bool("s", false,
		"Start a tunnel with the passed parameters(Otherwise, they will be treated as default values.)")
	encryptKeyFiles = flag.String("cr", "",
		"Encrypt/decrypt the key files with a passfile")
	inAllowZeroHop = flag.Bool("zi", false,
		"Allow zero-hop, non-anonymous tunnels in(true or false)")
	outAllowZeroHop = flag.Bool("zo", false,
		"Allow zero-hop, non-anonymous tunnels out(true o false)")
	useCompression = flag.Bool("z", false,
		"Uze gzip(true or false)")
	targetHost = flag.String("h", "127.0.0.1",
		"Target host(Host of service to forward to i2p)")
	targetPort = flag.String("p", "8081",
		"Target port(Port of service to forward to i2p)")
	reduceIdle = flag.Bool("r", true,
		"Reduce tunnel quantity when idle(true or false)")
	closeIdle = flag.Bool("x", true,
		"Close tunnel idle(true or false)")
	targetDir = flag.String("d", "./tunnels/",
		"Directory to save tunnel configuration file in.")
	targetDest = flag.String("de", "",
		"Destination to connect client's to by default.")
	iniFile = flag.String("f", "none",
		"Use an ini file for configuration(config file options override passed arguments for now.)")
	samHost = flag.String("sh", "127.0.0.1",
		"SAM host")
	samPort = flag.String("sp", "7656",
		"SAM port")
	tunName = flag.String("n", "socks",
		"Tunnel name, this must be unique but can be anything.")
	inLength = flag.Int("il", 3,
		"Set inbound tunnel length(0 to 7)")
	outLength = flag.Int("ol", 3,
		"Set outbound tunnel length(0 to 7)")
	inQuantity = flag.Int("iq", 2,
		"Set inbound tunnel quantity(0 to 15)")
	outQuantity = flag.Int("oq", 2,
		"Set outbound tunnel quantity(0 to 15)")
	inVariance = flag.Int("iv", 0,
		"Set inbound tunnel length variance(-7 to 7)")
	outVariance = flag.Int("ov", 0,
		"Set outbound tunnel length variance(-7 to 7)")
	inBackupQuantity = flag.Int("ib", 1,
		"Set inbound tunnel backup quantity(0 to 5)")
	outBackupQuantity = flag.Int("ob", 1,
		"Set outbound tunnel backup quantity(0 to 5)")
	reduceIdleTime = flag.Int("rt", 600000,
		"Reduce tunnel quantity after X (milliseconds)")
	closeIdleTime = flag.Int("ct", 600000,
		"Reduce tunnel quantity after X (milliseconds)")
	reduceIdleQuantity = flag.Int("rq", 1,
		"Reduce idle tunnel quantity to X (0 to 5)")
	readKeys = flag.String("conv", "", "Display the base32 and base64 values of a specified .i2pkeys file")
)

var (
	err    error
	config *i2ptunconf.Conf
)

func main() {
	flag.Parse()

	if *readKeys != "" {

	}

	config = i2ptunconf.NewI2PBlankTunConf()
	if *iniFile != "none" && *iniFile != "" {
		config, err = i2ptunconf.NewI2PTunConf(*iniFile)
	} else {
		*startUp = true
	}
	config.TargetHost = config.GetHost(*targetHost, "127.0.0.1")
	config.TargetPort = config.GetPort(*targetPort, "8081")
	config.SaveFile = config.GetSaveFile(true, true)
	config.SaveDirectory = config.GetDir(*targetDir, "./tunnels/")
	config.SamHost = config.GetSAMHost(*samHost, "127.0.0.1")
	config.SamPort = config.GetSAMPort(*samPort, "7656")
	config.TunName = config.GetKeys(*tunName, "socks")
	config.InLength = config.GetInLength(*inLength, 3)
	config.OutLength = config.GetOutLength(*outLength, 3)
	config.InVariance = config.GetInVariance(*inVariance, 0)
	config.OutVariance = config.GetOutVariance(*outVariance, 0)
	config.InQuantity = config.GetInQuantity(*inQuantity, 2)
	config.OutQuantity = config.GetOutQuantity(*outQuantity, 2)
	config.InBackupQuantity = config.GetInBackups(*inBackupQuantity, 1)
	config.OutBackupQuantity = config.GetOutBackups(*outBackupQuantity, 1)
	config.InAllowZeroHop = config.GetInAllowZeroHop(*inAllowZeroHop, false)
	config.OutAllowZeroHop = config.GetOutAllowZeroHop(*outAllowZeroHop, false)
	config.UseCompression = config.GetUseCompression(*useCompression, true)
	config.ReduceIdle = config.GetReduceOnIdle(*reduceIdle, true)
	config.ReduceIdleTime = config.GetReduceIdleTime(*reduceIdleTime, 600000)
	config.ReduceIdleQuantity = config.GetReduceIdleQuantity(*reduceIdleQuantity, 2)
	config.CloseIdle = config.GetCloseOnIdle(*closeIdle, true)
	config.CloseIdleTime = config.GetCloseIdleTime(*closeIdleTime, 600000)
	config.KeyFilePath = config.GetKeyFile(*encryptKeyFiles, "")
	config.ClientDest = config.GetClientDest(*targetDest, "", "")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	if tunsocks, tunerr := tunmanager.NewManager(config.SamHost, config.SamPort, config.SaveDirectory, config.Print()); tunerr == nil {
		go func() {
			for sig := range c {
				if sig == os.Interrupt {
					if err := tunsocks.Cleanup(); err != nil {
						log.Println(err.Error())
					}
				}
			}
		}()
		tunsocks.Serve()
	} else {
		panic(tunerr)
	}
}
