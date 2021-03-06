package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

var svc ucp.ServiceQuery

var prevSpec bool

func init() {
	// Service flags
	ucpServiceList.Flags().StringVar(&svc.ServiceName, "name", "", "Examine a service by name")

	// Query options
	ucpServiceList.Flags().BoolVar(&svc.ID, "id", false, "Display task ID")
	ucpServiceList.Flags().BoolVar(&svc.Networks, "networks", false, "Display task Network connections")
	ucpServiceList.Flags().BoolVar(&svc.State, "state", false, "Display task state")
	ucpServiceList.Flags().BoolVar(&svc.Node, "node", false, "Display Node running task")
	ucpServiceList.Flags().BoolVar(&svc.Resolve, "resolve", false, "Resolve Task IDs to human readable names")

	// Service Reap flags
	ucpServiceReap.Flags().StringVar(&svc.ServiceName, "name", "", "Examine a service by name")

	// Service Architecture flags
	ucpServiceArchitecture.Flags().BoolVar(&prevSpec, "previousSpec", false, "Display the previous Service specification")

	// Add Service to UCP root commands
	UCPRoot.AddCommand(ucpService)

	// Add reap to service subcommands
	ucpService.AddCommand(ucpServiceList)
	ucpService.AddCommand(ucpServiceReap)
	ucpService.AddCommand(ucpServiceArchitecture)
}

var ucpService = &cobra.Command{
	Use:   "service",
	Short: "Interact with services",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpServiceList = &cobra.Command{
	Use:   "list",
	Short: "List services",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		log.Debugf("Looking for service [%s]", svc.ServiceName)

		if svc.ServiceName != "" {
			err = client.QueryServiceContainers(&svc)
			if err != nil {
				log.Fatalf("%v", err)
			}
			return
		}

		err = client.GetAllServices()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpServiceReap = &cobra.Command{
	Use:   "reap",
	Short: "Clean a service (not implemented)",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		if svc.ServiceName != "" {
			err := client.QueryServiceContainers(&svc)
			if err != nil {
				log.Fatalf("%v", err)
			}
			return
		}
	},
}

var ucpServiceArchitecture = &cobra.Command{
	Use:   "architecture",
	Short: "Retrieve the \"design\" of a service",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		if len(args) == 0 {
			log.Fatalf("Please specify a service name")
		}
		if len(args) > 1 {
			log.Fatalf("Please only specify a single service to view the architecture")
		}

		log.Infof("Inspecting service [%s]", args[0])

		service, err := client.GetService(args[0])

		if err != nil {
			ucp.ParseUCPError([]byte(err.Error()))
			return
		}

		var spec *swarm.ServiceSpec

		if prevSpec == true {
			spec = service.PreviousSpec
		} else {
			spec = &service.Spec
		}

		printServiceSpec(service, spec)
	},
}

// This will read through the service spec and print out the details
func printServiceSpec(service *swarm.Service, spec *swarm.ServiceSpec) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	fmt.Fprintf(w, "ID:\t%s\n", service.ID)
	fmt.Fprintf(w, "Version:\t%d\n", service.Version.Index)
	fmt.Fprintf(w, "Name:\t%s\n", spec.Name)
	fmt.Fprintf(w, "Image:\t%s\n", spec.TaskTemplate.ContainerSpec.Image)
	//Print out the command used for the image
	fmt.Fprintf(w, "Cmd:")
	for i := range spec.TaskTemplate.ContainerSpec.Command {
		fmt.Fprintf(w, " %s", spec.TaskTemplate.ContainerSpec.Command[i])
	}
	fmt.Fprintf(w, "\n")
	// Print all arguments to the command
	fmt.Fprintf(w, "Args:")
	for i := range spec.TaskTemplate.ContainerSpec.Args {
		fmt.Fprintf(w, " %s", spec.TaskTemplate.ContainerSpec.Args[i])
	}
	fmt.Fprintf(w, "\n")

	// Print the labels from the key/map
	fmt.Fprintf(w, "Labels:\n")
	for key, value := range spec.TaskTemplate.ContainerSpec.Labels {
		fmt.Fprintf(w, "\t%s", key)
		fmt.Fprintf(w, "\t%s\n", value)
	}

	//Print reservations
	if spec.TaskTemplate.Resources != nil {
		fmt.Fprintf(w, "Memory Reservation:\t%d\n", spec.TaskTemplate.Resources.Reservations.MemoryBytes)
		fmt.Fprintf(w, "CPU Reservation:\t%d\n", spec.TaskTemplate.Resources.Reservations.NanoCPUs)
		//Print limits
		fmt.Fprintf(w, "Memory Limits:\t%d\n", spec.TaskTemplate.Resources.Limits.MemoryBytes)
		fmt.Fprintf(w, "CPU Limits:\t%d\n", spec.TaskTemplate.Resources.Limits.NanoCPUs)
	}
	// Check the pointer to the replica struct isn't nil and read replica count
	if spec.Mode.Replicated != nil {
		fmt.Fprintf(w, "Replicas:\t%d\n", *spec.Mode.Replicated.Replicas)
	}
	if spec.Mode.Global != nil {
		fmt.Fprintf(w, "Global:\ttrue\n")
	}
	w.Flush()

}

// {
// 	"ID": "qhg6qgcv6hm58fos2dl64ey43",
// 	"Version": {
// 	  "Index": 6794
// 	},
// 	"CreatedAt": "2018-05-16T18:24:06.556823733Z",
// 	"UpdatedAt": "2018-05-23T12:49:09.35299171Z",
// 	"Spec": {
// 	  "Name": "urchin",
// 	  "Labels": {
// 		"com.docker.ucp.access.label": "/",
// 		"com.docker.ucp.collection": "swarm",
// 		"com.docker.ucp.collection.root": "true",
// 		"com.docker.ucp.collection.swarm": "true"
// 	  },
// 	  "TaskTemplate": {
// 		"ContainerSpec": {
// 		  "Image": "thebsdbox/urchin:1.2@sha256:fbadb7d721cd9faabdead81323a02deb1a05993e3e60c0762eb249bed2d168d3",
// 		  "Labels": {
// 			"com.docker.ucp.access.label": "/",
// 			"com.docker.ucp.collection": "swarm",
// 			"com.docker.ucp.collection.root": "true",
// 			"com.docker.ucp.collection.swarm": "true"
// 		  },
// 		  "Command": [
// 			"/urchin"
// 		  ],
// 		  "Args": [
// 			"-w",
// 			"8080"
// 		  ],
// 		  "DNSConfig": {},
// 		  "Isolation": "default"
// 		},
// 		"Resources": {
// 		  "Limits": {
// 			"MemoryBytes": 8388608
// 		  },
// 		  "Reservations": {
// 			"MemoryBytes": 4194304
// 		  }
// 		},
// 		"Placement": {
// 		  "Constraints": [
// 			"node.labels.com.docker.ucp.collection.swarm==true",
// 			"node.labels.com.docker.ucp.orchestrator.swarm==true"
// 		  ],
// 		  "Platforms": [
// 			{
// 			  "Architecture": "amd64",
// 			  "OS": "linux"
// 			}
// 		  ]
// 		},
// 		"ForceUpdate": 0,
// 		"Runtime": "container"
// 	  },
// 	  "Mode": {
// 		"Replicated": {
// 		  "Replicas": 40
// 		}
// 	  },
// 	  "EndpointSpec": {
// 		"Mode": "vip"
// 	  }
// 	},
// 	"PreviousSpec": {
// 	  "Name": "urchin",
// 	  "Labels": {
// 		"com.docker.ucp.access.label": "/",
// 		"com.docker.ucp.collection": "swarm",
// 		"com.docker.ucp.collection.root": "true",
// 		"com.docker.ucp.collection.swarm": "true"
// 	  },
// 	  "TaskTemplate": {
// 		"ContainerSpec": {
// 		  "Image": "thebsdbox/urchin:1.2@sha256:fbadb7d721cd9faabdead81323a02deb1a05993e3e60c0762eb249bed2d168d3",
// 		  "Labels": {
// 			"com.docker.ucp.access.label": "/",
// 			"com.docker.ucp.collection": "swarm",
// 			"com.docker.ucp.collection.root": "true",
// 			"com.docker.ucp.collection.swarm": "true"
// 		  },
// 		  "Command": [
// 			"/urchin"
// 		  ],
// 		  "Args": [
// 			"-w",
// 			"8080"
// 		  ],
// 		  "DNSConfig": {},
// 		  "Isolation": "default"
// 		},
// 		"Resources": {
// 		  "Limits": {
// 			"MemoryBytes": 102410241
// 		  },
// 		  "Reservations": {}
// 		},
// 		"Placement": {
// 		  "Constraints": [
// 			"node.labels.com.docker.ucp.collection.swarm==true",
// 			"node.labels.com.docker.ucp.orchestrator.swarm==true"
// 		  ],
// 		  "Platforms": [
// 			{
// 			  "Architecture": "amd64",
// 			  "OS": "linux"
// 			}
// 		  ]
// 		},
// 		"ForceUpdate": 0,
// 		"Runtime": "container"
// 	  },
// 	  "Mode": {
// 		"Replicated": {
// 		  "Replicas": 40
// 		}
// 	  },
// 	  "EndpointSpec": {
// 		"Mode": "vip"
// 	  }
// 	},
// 	"Endpoint": {
// 	  "Spec": {}
// 	},
// 	"UpdateStatus": {
// 	  "State": "updating",
// 	  "StartedAt": "2018-05-23T12:49:09.352973272Z",
// 	  "Message": "update in progress"
// 	}
//   }
