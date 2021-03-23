package handlers

import (
	"context"
	"time"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/profiles"

	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"

	"bitbucket.org/edgelabsolutions/loncherapp-profiles-service/app/services"
)

type ProfilesAPIServer struct {
	ProfileService *services.ProfileService
}

func NewProfilesAPIServer(ctx context.Context) *ProfilesAPIServer {
	return &ProfilesAPIServer{
		ProfileService: services.NewProfileService(ctx),
	}
}

func (p *ProfilesAPIServer) DeleteProfileUserFavorite(ctx context.Context, request *pbm.ProfileUserFavoriteRequest) (*pbm.SimpleResponse, error) {
	err := p.ProfileService.DeleteProfileUserFavorite(request.UserId, request.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: true, Message: "delete favorite correctly"}, nil
}

func (p *ProfilesAPIServer) SetProfileUserFavorite(ctx context.Context, request *pbm.ProfileUserFavoriteRequest) (*pbm.SimpleResponse, error) {
	err := p.ProfileService.SetProfileUserFavorite(request.UserId, request.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: true, Message: "Set favorite correctly"}, nil
}

func (p *ProfilesAPIServer) CreateProfile(ctx context.Context, profile *pbm.Profile) (*pbm.SimpleResponse, error) {
	_, err := p.ProfileService.CreateProfile(ConvertProfilePMtoModel(profile))
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: true, Message: "Profile created correctly"}, nil
}

func (p *ProfilesAPIServer) UpdateProfile(ctx context.Context, profile *pbm.Profile) (*pbm.SimpleResponse, error) {
	_, err := p.ProfileService.UpdateProfile(ConvertProfilePMtoModel(profile))
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: true, Message: "Profile created correctly"}, nil
}

func (p *ProfilesAPIServer) GetProfileByUserID(ctx context.Context, simpleRequestByID *pbm.SimpleRequestByID) (*pbm.Profile, error) {
	response, err := p.ProfileService.GetProfileByUserID(simpleRequestByID.Id)
	if err != nil {
		return nil, err
	}

	return ConvertProfileModeltoPM(response), nil
}

func (p *ProfilesAPIServer) GetProfileByID(ctx context.Context, simpleRequestByID *pbm.SimpleRequestByID) (*pbm.Profile, error) {
	response, err := p.ProfileService.GetProfileByID(simpleRequestByID.Id)
	if err != nil {
		return nil, err
	}

	return ConvertProfileModeltoPM(response), nil
}

func (p *ProfilesAPIServer) GetProfilesByCategory(ctx context.Context, profileQueryCategory *pbm.ProfileQueryCategory) (*pbm.Profiles, error) {
	response, err := p.ProfileService.GetProfilesByCategory(profileQueryCategory.Category)
	if err != nil {
		return nil, err
	}

	var listPbmProfile pbm.Profiles

	for _, profile := range *response {
		listPbmProfile.Profiles = append(listPbmProfile.Profiles, ConvertProfileModeltoPM(&profile))
	}

	return &listPbmProfile, nil
}

func (p *ProfilesAPIServer) GetProfilesByUserLocation(ctx context.Context, profileQueryUserLocation *pbm.ProfileQueryUserLocation) (*pbm.Profiles, error) {

	panic("implement me")
}

func (p *ProfilesAPIServer) GetProfilesByUserFavorites(ctx context.Context, simpleRequestID *pbm.SimpleRequestByID) (*pbm.Profiles, error) {
	response, err := p.ProfileService.GetProfilesByUserFavorites(simpleRequestID.Id)
	if err != nil {
		return nil, err
	}

	var listPbmProfile pbm.Profiles

	for _, profile := range *response {
		modelProfile := ConvertProfileModeltoPM(&profile)
		modelProfile.IsFavorite = true
		listPbmProfile.Profiles = append(listPbmProfile.Profiles, modelProfile)
	}

	return &listPbmProfile, nil
}

func ConvertProfilePMtoModel(pbmProfile *pbm.Profile) *models.Profile {
	membershipDate, _ := time.Parse(time.RFC3339, pbmProfile.MembershipDueDate)
	createdAt, _ := time.Parse(time.RFC3339, pbmProfile.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, pbmProfile.UpdatedAt)

	return &models.Profile{
		ID:                pbmProfile.Id,
		UserID:            pbmProfile.UserId,
		Description:       pbmProfile.Description,
		CategoryID:        pbmProfile.CategoryId,
		CoverImageURL:     pbmProfile.CoverImageUrl,
		Website:           pbmProfile.Website,
		Active:            pbmProfile.Active,
		MembershipDueDate: membershipDate,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		Rating:            pbmProfile.Rating,
		IsFavorite:        pbmProfile.IsFavorite,
	}
}

func ConvertProfileModeltoPM(profile *models.Profile) *pbm.Profile {
	return &pbm.Profile{
		Id:                profile.ID,
		UserId:            profile.UserID,
		Description:       profile.Description,
		CategoryId:        profile.CategoryID,
		CoverImageUrl:     profile.CoverImageURL,
		Website:           profile.Website,
		Active:            profile.Active,
		MembershipDueDate: profile.MembershipDueDate.String(),
		CreatedAt:         profile.CreatedAt.String(),
		UpdatedAt:         profile.UpdatedAt.String(),
		Rating:            profile.Rating,
		IsFavorite:        profile.IsFavorite,
	}
}
