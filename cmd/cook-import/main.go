package main

import (
	"os"
	"regexp"
	"bytes"
	"path"
	"context"
	"net/url"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	openai "github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
)

const defaultMessage = `Here are the specification for a markup language called cooklang, used to describe a cooking recipe:
Define an ingredient using the @ symbol, indicate the end of the ingredient's name with {}, indicate the quantity of an item inside {} after the name,
indicate the unit of an item's quantity, such as weight or volume, using %, define any necessary cookware with #, define a timer using ~
Here is an example:
` +
	"```" +
	`
Crack the @eggs{3} into a blender, then add the @flour{125%g), @milk{250%ml} and @sea salt{1%pinch}, and blitz until smooth.

Pour into a #bowl and leave to stand for ~{15%minutes}.

Melt the @butter{} (or a drizzle of @oil if you want to be a bit healthier) in a #large non-stick frying pan{} on a medium heat, then tilt the pan so the butter coats the surface.

Pour in 1 ladle of batter and tilt again, so that the batter spreads all over the base, then cook for 1 to 2 minutes, or until it starts to come away from the sides.

Once golden underneath, flip the pancake over and cook for 1 further minute, or until cooked through.

Serve straightaway with your favorite topping. - Add your favorite topping here to make sure it's included in your meal plan!
` +
	"```" +
	`
Never make the list of ingredients, only write the instructions using the markup

Abbreviate ingredient measurements and convert fractions to decimals in ingredients

Using cooklang, get the recipe instructions at the link` // https://healthyrecipesblogs.com/crustless-quiche/

func getContent(servings string, instructions string, link string)(string) {
	content :=`>> source: ` + link + "\n" +
`>> serves: ` + servings +
"\n\n" + instructions + "\n"
	return content
}

func getOutputFile(newFileName string, toFile bool) (*os.File, error) {
	if toFile {
		return os.Create(newFileName)
	}
	return os.Stdout, nil
}

func applyMarkDownFormat(output bytes.Buffer) bytes.Buffer {
	outputString := output.String()
	re := regexp.MustCompile(` \n`)
	outputString = re.ReplaceAllString(outputString, "\n")

	re = regexp.MustCompile(`\n{3,}`)
	outputString = re.ReplaceAllString(outputString, "\n\n")

	output.Reset()
	output.WriteString(outputString)
	return output
}

func getFilePath(title string) (string, error){
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	filepath := path.Join(wd,title + ".cook")
	return filepath, nil
}

func printDocumentation(title string, content string) {
	toFile := viper.GetBool("file")

	filepath, err := getFilePath(title)
	if err != nil {
		log.Errorf("Could not get filepath %v", err)
		return
	}
	
	outputFile, err := getOutputFile(filepath, toFile)
	if err != nil {
			log.Errorf("Could not get the output file %v\n", err)
			return
	}	
	if toFile{
		defer outputFile.Close()
		log.Infof("Writing to file %v", filepath)
	}
	var output bytes.Buffer
	output.WriteString(content)
	output = applyMarkDownFormat(output)
	_, err = output.WriteTo(outputFile)
	if err != nil {
		log.Errorf("Error generating cook file for recipe %v: %v", filepath, err)
	}
}

func cookImport(_ *cobra.Command, _ []string) {
	initializeCli()
	link := viper.GetString("link")
	_, err := url.ParseRequestURI(link)

	if err != nil {
		log.Errorf("Invalid URI: %v\n", err)
		return
	}
	log.Debugf("link: %v\n", link)

	message := defaultMessage + " " + link
	log.Debugf("message: %v\n", message)
	
	key := viper.GetString("openai-api-key")
	if len(key) == 0 {
		log.Errorln("OpenAI API key is not set")
		return
	}
	log.Debugf("key: %v\n", key)

	client := openai.NewClient(key)
	ctx := context.Background()
	model := openai.GPT3Dot5Turbo
	
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		},
	}
	log.Infof("Getting instructions from %v\n", link)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)
	if err != nil {
		log.Errorf("ChatCompletion error: %v\n", err)
		return
	}
	instructions := resp.Choices[0].Message.Content
	log.Debugf("instructions: %v\n", instructions)
	messages = []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Get the number of servings of the recipe at this link " + link + ". Return only the digit",
		},
	}
	log.Infoln("Getting servings")
	resp, err = client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)
	if err != nil {
		log.Errorf("ChatCompletion error: %v\n", err)
		return
	}
	re := regexp.MustCompile("[0-9]+")
	servings := re.FindString(resp.Choices[0].Message.Content)
	log.Infof("servings: %v\n", servings)

	messages = []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Gget the title of the recipe from the link " + link,
		},
	}
	log.Infoln("Getting title")
	resp, err = client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)
	if err != nil {
		log.Errorf("ChatCompletion error: %v\n", err)
		return
	}
	title := resp.Choices[0].Message.Content
	log.Infof("title: %v\n", title)
	
	log.Infoln("Generating recipe")
	content:=getContent(servings,instructions,link)
	log.Debugf("content: %v\n", content)
	
	printDocumentation(title, content)
}

func main() {
	command, err := newCookImportCommand(cookImport)
	if err != nil {
		log.Errorf("Failed to create the CLI command: %v", err)
		os.Exit(1)
	}

	if err := command.Execute(); err != nil {
		log.Errorf("Failed to start the CLI: %v", err)
		os.Exit(1)
	}
}
