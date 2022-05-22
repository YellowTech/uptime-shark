package notifications

import (
	"context"
	"log"
	"strconv"
	"uptime/api/ent"
)

type NotificationService struct {
    Name string
	SendMessage func(string) error
}

func (n NotificationService) String() string {
        return n.Name
}

// load services from DB and setting them up
func SetupNotifications(client *ent.Client, ctx context.Context) ([]*NotificationService, error) {
	notificationEntries, err := client.Notification.Query().All(ctx)

	if (err != nil) {
		log.Println("Error occurred during loading of notification entries")
		return nil, err
	}

	services := make([]*NotificationService, 0)

	for _, entry := range notificationEntries {
		// if not active, skip
		if (!entry.Active){
			continue
		}
		var newService *NotificationService
		
		switch entry.Settings[0] {
		case "telegram":
			if(len(entry.Settings) != 3) {
				log.Println("Error loading notification service of type " + entry.Name + ", wrong amount of saved settings")
				continue
			}

			apiKey := entry.Settings[1]
			var chatId int64
			chatId, err = strconv.ParseInt(entry.Settings[2], 10, 64)
			if (err != nil) {
				log.Println("Error loading notification service of type " + entry.Name, err)
				continue
			}

			newService, err = CreateTelegramService(apiKey, chatId, entry.Name)
			if (err != nil) {
				log.Println("Error creating telegram service", err)
				continue
			}

		default:
			log.Println("Error loading notification service of type " + entry.Name)
			continue
		}
		newService.SendMessage("Initialized Notification Service")
		services = append(services, newService)		
	}

    return services, err
}

func SendMessageToAll(services []*NotificationService, message string) {
	for _, service := range services {
		service.SendMessage(message);
	}
}