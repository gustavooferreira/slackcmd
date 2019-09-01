package security

type Permissions struct {
	TeamID  string
	Global  []string // If GO had SETs, that would be best.
	Channel map[string][]string
}

func NewPermissions(teamID string, globalPermissions []string, channelPermissions map[string][]string) Permissions {
	if channelPermissions == nil {
		channelPermissions = make(map[string][]string)
	}
	return Permissions{TeamID: teamID, Global: globalPermissions, Channel: channelPermissions}
}

func (p *Permissions) AddGlobalUser(userID string) {
	p.Global = append(p.Global, userID)
}

func (p *Permissions) AddChannelUser(channelID string, userID string) {
	if value, ok := p.Channel[channelID]; ok {
		value = append(value, userID)
	} else {
		p.Channel[channelID] = []string{userID}
	}
}

func (p Permissions) ValidateGlobal(userID string) bool {
	for _, elem := range p.Global {
		if elem == userID {
			return true
		}
	}
	return false
}

func (p Permissions) ValidateChannel(channelID string, userID string) bool {

	value, ok := p.Channel[channelID]
	if !ok {
		return false
	}

	for _, elem := range value {
		if elem == userID {
			return true
		}
	}
	return false
}
