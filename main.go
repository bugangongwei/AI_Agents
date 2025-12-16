package main

import (
	"fmt"
	"os"
	"os/exec"
	outfit_recommender "outfit-recommender/outfit-recommender"
)

type Agent struct {
	Name        string
	Description string
	Usage       string
	Path        string
}

var agents = map[string]Agent{
	"outfit-recommender": {
		Name:        "outfit-recommender",
		Description: "Recommends daily outfits based on weather, schedule, user preferences, and input using AI",
		Usage:       "go run . outfit-recommender -input \"<user description>\" -pref <style> -loc \"<city>\"\n  -input: Text describing your preferences (e.g., \"I like blue colors\")\n  -pref: Style preference (casual, formal, sporty, elegant)\n  -loc: City name for weather data (e.g., \"Beijing\")",
		Path:        "outfit-recommender/start.go",
	},
}

func main() {
	outfit_recommender.LoadFashionTrend()

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	if command == "--help" || command == "-h" {
		printHelp()
		return
	}

	if agent, ok := agents[command]; ok {
		runAgent(agent, os.Args[2:])
		return
	}

	fmt.Printf("Unknown agent: %s\n", command)
	printHelp()
}

func printHelp() {
	fmt.Println("AI Agents Collection")
	fmt.Println("Usage: go run . <agent-name> [args...]")
	fmt.Println()
	fmt.Println("Available agents:")
	for name, agent := range agents {
		fmt.Printf("	%s: %s\n", name, agent.Description)
		fmt.Printf("	Usage: %s\n", agent.Usage)
		fmt.Println()
	}
	fmt.Println("Use 'go run . --help' for this message")
}

func runAgent(agent Agent, args []string) {
	if agent.Name == "outfit-recommender" {
		outfit_recommender.Start(args)
	} else {
		cmd := exec.Command("go", "run", agent.Path)
		cmd.Args = append(cmd.Args, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error running agent %s: %v\n", agent.Name, err)
			os.Exit(1)
		}
	}
}
