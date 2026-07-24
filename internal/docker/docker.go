package docker

import (
	"fmt"
	"log"
	"os/exec"
)

func Run(image string) error {

	log.Println("[DOCKER] Preparing docker command")

	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		image,
	)

	log.Printf("[DOCKER] Executing: docker run -d %s\n", image)

	output, err := cmd.CombinedOutput()

	log.Printf("[DOCKER] Output:\n%s\n", string(output))

	if err != nil {
		log.Printf("[DOCKER] Docker command failed")
		return fmt.Errorf("%v\n%s", err, output)
	}

	log.Println("[DOCKER] Docker command completed successfully")

	return nil
}