package bc

import "time"

type bcApiUsersList struct {
	Status  string
	Message string
	Object  []*bcApiUser
}

type bcApiUser struct {
	Id            string
	ServiceUserId string `json:"username"`
	Services      map[string]interface{}
	Metadata      map[string]string
	Customer      string
}

func toUserResource(u *bcApiUser) *User {
	user := &User{
		Id: u.ServiceUserId,
		// Id:     u.Id,
		Email:  u.Metadata["facebook.user.email"],
		Gender: u.Metadata["facebook.user.gender"],
		About:  u.Metadata["twitter.user.description"],
	}

	user.Name = u.Metadata["facebook.user.name"]
	if user.Name == "" {
		user.Name = u.Metadata["twitter.user.name"]
	}

	user.Location = u.Metadata["facebook.user.locationname"]
	if user.Location == "" {
		user.Location = u.Metadata["twitter.user.location"]
	}

	// mm/dd/YYYY
	dob, _ := time.Parse("01/02/2006", u.Metadata["facebook.user.birthday"])
	if dob.IsZero() {
		user.Age = 0
	} else {
		user.Age = int(time.Since(dob) / (time.Hour * 24 * 30 * 12))
	}

	if _, ok := u.Services["twitter"]; ok {
		user.Services = append(user.Services, &UserService{
			Id:    u.Metadata["twitter.user.screenName"],
			Name:  "twitter",
			Link:  "https://twitter.com/" + u.Metadata["twitter.user.screenName"],
			Photo: u.Metadata["twitter.user.imageUrl"],
		})
	}

	if _, ok := u.Services["facebook"]; ok {
		user.Services = append(user.Services, &UserService{
			Id:    u.Metadata["facebook.user.id"],
			Name:  "facebook",
			Link:  u.Metadata["facebook.user.link"],
			Photo: u.Metadata["facebook.user.picture"],
		})
	}

	for _, s := range user.Services {
		if s.Photo != "" {
			user.Photo = s.Photo
			break
		}
	}

	return user
}

type bcApiUserProfile struct {
	Status  string
	Message string
	Object  struct {
		LastUpdated int64
		Interests   []struct {
			Resource   string
			Label      string
			Weight     float32
			Activities []string
		}
		Categories []struct {
			Resource string
			Label    string
			Weight   interface{}
			Urls     []string
		}
	}
}

func toUserProfileResource(p *bcApiUserProfile) *UserProfile {
	timestamp := p.Object.LastUpdated / 1000
	profile := &UserProfile{
		Updated: time.Unix(timestamp, (p.Object.LastUpdated-timestamp)*1000),
		Topics:  make([]*Topic, 0, len(p.Object.Interests)+len(p.Object.Categories)),
	}

	for _, interest := range p.Object.Interests {
		profile.Topics = append(profile.Topics, &Topic{
			Kind:       "topic",
			Resource:   interest.Resource,
			Label:      interest.Label,
			Weight:     interest.Weight,
			Activities: interest.Activities,
		})
	}

	for _, cat := range p.Object.Categories {
		topic := &Topic{
			Kind:     "interest",
			Resource: cat.Resource,
			Label:    cat.Label,
			Urls:     cat.Urls,
		}
		if w, ok := cat.Weight.(float32); ok {
			topic.Weight = w
		}
		profile.Topics = append(profile.Topics, topic)
	}

	return profile
}

type bcApiActivitiesList struct {
	Status  string
	Message string
	Object  []struct {
		ApiActivity bcApiActivity `json:"activity"`
	}
}

type bcApiActivity struct {
	Id     string
	Verb   string
	Object struct {
		Type        string
		Url         string
		Name        string
		Description string
	}
	Context struct {
		Date    int64
		Service string
	}
}

func toActivityResource(a *bcApiActivity) *Activity {
	timestamp := a.Context.Date / 1000
	return &Activity{
		Id:          a.Id,
		Verb:        a.Verb,
		Kind:        a.Object.Type,
		Service:     a.Context.Service,
		Url:         a.Object.Url,
		Name:        a.Object.Name,
		Description: a.Object.Description,
		Timestamp:   time.Unix(timestamp, (a.Context.Date-timestamp)*1000),
	}
}
