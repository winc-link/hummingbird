package container

import (
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/tools/agentclient"
)

var AgentClientName = di.TypeInstanceToName((*agentclient.AgentClient)(nil))

func AgentClientNameFrom(get di.Get) agentclient.AgentClient {
	return get(AgentClientName).(agentclient.AgentClient)
}
