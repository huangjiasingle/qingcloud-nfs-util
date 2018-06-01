package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	qcfg "github.com/yunify/qingcloud-sdk-go/config"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

var accessKeyID = flag.String("accessKeyID", "", "")
var secretAccessKey = flag.String("secretAccessKey", "", "")
var host = flag.String("host", "", "")
var port = flag.Int("port", 0, "")
var zone = flag.String("zone", "", "")

func validateFlags() {
	flag.Parse()
	if *accessKeyID == "" {
		if os.Getenv("ACCESS_KEY_ID") != "" {
			*accessKeyID = os.Getenv("ACCESS_KEY_ID")
		} else {
			log.Fatal("accessKeyID flag or env ACCESS_KEY_IDmust define")
		}
	}
	if *secretAccessKey == "" {
		if os.Getenv("SECRET_ACCESS_KEY") != "" {
			*secretAccessKey = os.Getenv("SECRET_ACCESS_KEY")
		} else {
			log.Fatal("secretAccessKey flag or env SECRET_ACCESS_KEY define")
		}
	}
	if *host == "" {
		if os.Getenv("HOST") != "" {
			*host = os.Getenv("HOST")
		} else {
			log.Fatal("host flag or env HOST define")
		}
	}
	if *port == 0 {
		if os.Getenv("PORT") != "" {
			p, _ := strconv.Atoi(os.Getenv("PORT"))
			*port = p
		} else {
			log.Fatal("port flag or env HOST define")
		}
	}
	if *zone == "" && os.Getenv("ZONE") == "" {
		if os.Getenv("PORT") != "" {
			*zone = os.Getenv("ZONE")
		} else {
			log.Fatal("zone flag or env ZONE define")
		}
	}
}

func initQcClient() (*qc.QingCloudService, error) {
	configuration, _ := qcfg.New(*accessKeyID, *secretAccessKey)
	configuration.Protocol = "http"
	configuration.Host = *host
	configuration.Port = *port
	qcService, err := qc.Init(configuration)
	return qcService, err
}

func main() {
	validateFlags()
	qclient, err := initQcClient()
	if err != nil {
		log.Fatal(err)
	}

	volumeService, err := qclient.Volume(*zone)
	if err != nil {
		log.Fatal(err)
	}

	queryParam := &qc.DescribeVolumesInput{
		Limit:  qc.Int(10000),
		Offset: qc.Int(0),
		Status: []*string{qc.String("available")},
	}
	volumeList, err := volumeService.DescribeVolumes(queryParam)
	if err != nil {
		log.Fatal(err)
	}

	deleteParam := &qc.DeleteVolumesInput{}
	ids := []*string{}
	for _, vol := range volumeList.VolumeSet {
		if strings.Contains(*vol.VolumeName, "delete") {
			ids = append(ids, vol.VolumeID)
		}
	}
	if len(ids) == 0 {
		log.Println("no nfs's volume need to delete")
		return
	}
	deleteParam.Volumes = ids
	if _, err := volumeService.DeleteVolumes(deleteParam); err != nil {
		log.Printf("delete qingcloud nfs's volume err: %v", err)
	} else {
		log.Print("delete qingcloud nfs's volume whhere need to delete success")
	}
}
