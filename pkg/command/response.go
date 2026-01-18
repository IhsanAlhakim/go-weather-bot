package command

import "github.com/bwmarrin/discordgo"

func singleTextResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message *string) {
	// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Content: message,
	// 	},
	// })

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: message,
	})
}

func deferInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
		return err
	}
	return nil
}
