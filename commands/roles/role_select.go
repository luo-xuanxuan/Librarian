package roles

import (
	"Librarian/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func create_role_select(category *role_category, i *discordgo.InteractionCreate) *discordgo.SelectMenu {
	var roles map[string]string = category.Roles
	var role_options = make([]discordgo.SelectMenuOption, 0)

	for name, id := range roles {
		utils.Log.WithField("name", id).Info()
		role_options = append(role_options, discordgo.SelectMenuOption{
			Label:       name,
			Value:       id,
			Description: "",
			Default:     utils.Contains(i.Member.Roles, id),
		})
	}

	max_value := category.MaxValues
	if max_value > len(role_options) {
		max_value = len(role_options)
	}

	return &discordgo.SelectMenu{
		CustomID:    "role_select",
		Placeholder: "Select a Role.",
		Options:     role_options,
		MinValues:   &category.MinValues,
		MaxValues:   max_value,
	}
}

// role select handler
func role_select_handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	component := i.MessageComponentData()

	if !strings.HasPrefix(component.CustomID, "role_select") {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		fmt.Println("Failed to respond to interaction:", err)
		return
	}

	parts := strings.Split(component.CustomID, "_")

	// Retrieve the last item
	index, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		utils.Log.Error(err)
		return
	}

	var category *role_category = &cfg[i.GuildID][index]

	//available roles
	var role_list []string
	for _, id := range category.Roles {
		role_list = append(role_list, id)
	}

	// Calculate roles to add and remove directly
	roles_to_add := utils.Difference(component.Values, i.Member.Roles)
	roles_to_remove := utils.Difference(role_list, component.Values)
	roles_to_remove = utils.Intersection(roles_to_remove, i.Member.Roles)

	guild_roles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		utils.Log.Error(err)
	}

	roles_map := make(map[string]string)
	for _, gr := range guild_roles {
		roles_map[gr.ID] = gr.Name
	}

	// Assign new roles
	for j, role := range roles_to_add {
		utils.Assign_Role(s, i.GuildID, i.Member.User.ID, role)
		roles_to_add[j] = roles_map[role]
	}
	utils.Log.WithField("User", i.Member.User.Username).
		WithField("Guild", i.GuildID).
		WithField("Roles", strings.Join(roles_to_add, ", ")).
		Info("Added Roles")

	// Remove unselected roles
	for j, role := range roles_to_remove {
		utils.Remove_Role(s, i.GuildID, i.Member.User.ID, role)
		roles_to_remove[j] = roles_map[role]
	}
	utils.Log.WithField("User", i.Member.User.Username).
		WithField("Guild", i.GuildID).
		WithField("Roles", strings.Join(roles_to_remove, ", ")).
		Info("Removed Roles")

}
