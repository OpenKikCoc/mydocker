package container

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//Create a AUFS filesystem as container root workspace
func NewWorkSpace(volume, imageName, containerName string) {
	CreateReadOnlyLayer(imageName)
	CreateWriteLayer(containerName)
	CreateMountPoint(containerName, imageName)
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(volumeURLs, containerName)
			log.Printf("NewWorkSpace volume urls %q\n", volumeURLs)
		} else {
			log.Printf("Volume parameter input is not correct.\n")
		}
	}
}

//Decompression tar image
func CreateReadOnlyLayer(imageName string) error {
	unTarFolderUrl := RootUrl + "/" + imageName + "/"
	imageUrl := RootUrl + "/" + imageName + ".tar"
	exist, err := PathExists(unTarFolderUrl)
	if err != nil {
		log.Printf("Fail to judge whether dir %s exists. %v\n", unTarFolderUrl, err)
		return err
	}
	if !exist {
		if err := os.MkdirAll(unTarFolderUrl, 0622); err != nil {
			log.Printf("Mkdir %s error %v\n", unTarFolderUrl, err)
			return err
		}

		if _, err := exec.Command("tar", "-xvf", imageUrl, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			log.Printf("Untar dir %s error %v\n", unTarFolderUrl, err)
			return err
		}
	}
	return nil
}

func CreateWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		log.Printf("Mkdir write layer dir %s error. %v\n", writeURL, err)
	}
}

func MountVolume(volumeURLs []string, containerName string) error {
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil {
		log.Printf("Mkdir parent dir %s error. %v\n", parentUrl, err)
	}
	containerUrl := volumeURLs[1]
	mntURL := fmt.Sprintf(MntUrl, containerName)
	containerVolumeURL := mntURL + "/" + containerUrl
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
		log.Printf("Mkdir container dir %s error. %v\n", containerVolumeURL, err)
	}
	dirs := "dirs=" + parentUrl
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL).CombinedOutput()
	if err != nil {
		log.Printf("Mount volume failed. %v\n", err)
		return err
	}
	return nil
}

func CreateMountPoint(containerName, imageName string) error {
	mntUrl := fmt.Sprintf(MntUrl, containerName)
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		log.Printf("Mkdir mountpoint dir %s error. %v\n", mntUrl, err)
		return err
	}
	tmpWriteLayer := fmt.Sprintf(WriteLayerUrl, containerName)
	tmpImageLocation := RootUrl + "/" + imageName
	mntURL := fmt.Sprintf(MntUrl, containerName)
	dirs := "dirs=" + tmpWriteLayer + ":" + tmpImageLocation
	log.Println("CreateMountPoint RootUrl", RootUrl, "imageName", imageName)
	log.Println("CreateMountPoint cmd:", "mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	// mount -t aufs -o dirs=/root/writeLayer/9778725254:/root//bin/sh none /root/mnt/9778725254
	// root后面有两个// 问题出在tmpImageLocation
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL).CombinedOutput()
	if err != nil {
		log.Printf("Run command for creating mount point failed %v\n", err)
		return err
	}
	return nil
}

//Delete the AUFS filesystem while container exit
func DeleteWorkSpace(volume, containerName string) {
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteVolume(volumeURLs, containerName)
		}
	}
	DeleteMountPoint(containerName)
	DeleteWriteLayer(containerName)
}

func DeleteMountPoint(containerName string) error {
	mntURL := fmt.Sprintf(MntUrl, containerName)
	_, err := exec.Command("umount", mntURL).CombinedOutput()
	if err != nil {
		log.Printf("Unmount %s error %v\n", mntURL, err)
		return err
	}
	if err := os.RemoveAll(mntURL); err != nil {
		log.Printf("Remove mountpoint dir %s error %v\n", mntURL, err)
		return err
	}
	return nil
}

func DeleteVolume(volumeURLs []string, containerName string) error {
	mntURL := fmt.Sprintf(MntUrl, containerName)
	containerUrl := mntURL + "/" + volumeURLs[1]
	if _, err := exec.Command("umount", containerUrl).CombinedOutput(); err != nil {
		log.Printf("Umount volume %s failed. %v\n", containerUrl, err)
		return err
	}
	return nil
}

func DeleteWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		log.Printf("Remove writeLayer dir %s error %v\n", writeURL, err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
