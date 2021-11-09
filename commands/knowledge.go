package commands

import (
	"errors"
	"fmt"
	_ "github.com/jfrog/jfrog-cli-core/v2/common/commands"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-plugin-template/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func init() {}

func GetKnowledgeCommand() components.Command {
	return components.Command{
		Name:        "knowledge",
		Description: "Query JFrog Knowledge Base.",
		Aliases:     []string{"know"},
		Action: func(c *components.Context) error {
			return knowledgeCmd(c)
		},
	}
}

type knowledgeConfiguration struct {
	query string
}

func knowledgeCmd(c *components.Context) error {
	if len(c.Arguments) != 1 {
		return errors.New("No search expression provided. Please use the following format: 'jfrog know <query>'")
	}
	var conf = new(knowledgeConfiguration)
	conf.query = c.Arguments[0]

	facet, err := getFacetForQuery(conf.query)
	if err != nil {
		return err
	}

	postId, err := getResultsForQuery(facet, conf.query)
	if err != nil {
		return err
	}

	err = getContentForQuery(postId, conf.query)
	if err != nil {
		return err
	}

	return nil
}

func getFacetForQuery(query string) (string, error) {
	log.Info(fmt.Sprintf("Searching for '%s' in JFrog Knowledge Base...", query))

	facets, err := utils.GetFacetsContent("/api/v1/search/facets?query='" + url.QueryEscape(query) + "'")

	if err != nil {
		log.Warn(fmt.Sprintf("Could not find any search results for query: '%s'", query))
		os.Exit(0)
	}

	promptItems := []utils.PromptItem{}

	for facetName, facetQty := range facets {
		qty := strconv.Itoa(facetQty)
		promptItems = append(promptItems, utils.PromptItem{Id: 0, Option: facetName, TargetValue: &qty})
	}

	option, err := utils.PromptStringsNew(promptItems, "Select Content Type Filter:")

	facet := option.Option

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	return facet, nil
}

func getResultsForQuery(facet string, query string) (int64, error) {
	log.Info(fmt.Sprintf("Searching for '%s' with content type filter '%s' in JFrog Knowledge Base...", query, facet))

	var results []utils.KnowResult

	results, err := utils.GetResultsContent("/api/v1/search?query=" + url.QueryEscape(query) + "&facet=" + url.QueryEscape(facet))

	promptItems := []utils.PromptItem{}

	for _, res := range results {
		title := res.Title
		content := res.PublishDate + " [Author: " + res.Author + "]"
		promptItems = append(promptItems, utils.PromptItem{Id: res.PostID, Option: title, TargetValue: &content})
	}

	res, err := utils.PromptStringsNew(promptItems, fmt.Sprintf("Select %s:", facet))

	if err != nil {
		fmt.Println(err)
	}

	postId := res.Id

	return postId, nil

}

func getContentForQuery(postId int64, query string) error {
	log.Info(fmt.Sprintf("Searching for '%s' in JFrog Knowledge Base...", query))

	content, err := utils.GetContent("/api/v1/search/id/" + strconv.FormatInt(postId, 10))

	promptItems := []utils.PromptItem{
		{Id: postId, Option: "Open URL in default browser", TargetValue: &content.URL},
		{Id: postId, Option: "Read in the Terminal"},
	}

	choice, err := utils.PromptStringsNew(promptItems, fmt.Sprintf("How do you want to access content [ID: %d]?", postId))

	if err != nil {
		fmt.Println(err)
	}

	if choice.Option == "Open URL in default browser" {
		openBrowser(content.URL)
	} else {
		fmt.Sprintf("%s | Author: %s | Publish Date: %s", content.Title, content.Author, content.PublishDate)
		fmt.Println("")
		fmt.Println(fmt.Sprintf("%s ...", content.Content))
		fmt.Println("")
		fmt.Println(fmt.Sprintf("For the complete content, access the page here: %s", content.URL))
	}

	return nil
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
	}
}
