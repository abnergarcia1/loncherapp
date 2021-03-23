package handlers

import (
	"context"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/menus"

	pbm "bitbucket.org/edgelabsolutions/loncherapp-protobuf/go_proto/models"

	"bitbucket.org/edgelabsolutions/loncherapp-menus-service/app/services"
)

//MenusAPIServer Routes API Server
type MenusAPIServer struct {
	MenusService *services.MenusService
}

//GetMenuByProfileID Get Menu object list from DB using ProfileID
func (r MenusAPIServer) GetMenuByProfileID(ctx context.Context, id *pbm.SimpleRequestByID) (*pbm.Menus, error) {
	menus, err := r.MenusService.GetMenuByProfileID(id.Id)
	if err != nil {
		return nil, err
	}

	var listMenu = make([]*pbm.Menu, len(*menus))
	for i, menu := range *menus {
		listMenu[i] = ConvertMenuModelToPM(&menu)
	}

	return &pbm.Menus{Menus: listMenu}, nil
}

//GetMenuByID Get Menu element object from DB using MenuID
func (r MenusAPIServer) GetMenuByID(ctx context.Context, id *pbm.SimpleRequestByID) (*pbm.Menu, error) {
	menu, err := r.MenusService.GetMenuByID(id.Id)
	if err != nil {
		return nil, err
	}

	return ConvertMenuModelToPM(menu), nil
}

//CreateMenu Create new Menu element object in DB
func (r MenusAPIServer) CreateMenu(ctx context.Context, menu *pbm.Menu) (*pbm.Menu, error) {
	modelMenu := ConvertMenuPMtoModel(menu)
	menuResponse, err := r.MenusService.CreateMenu(modelMenu)
	if err != nil {
		return nil, err
	}

	return ConvertMenuModelToPM(menuResponse), nil
}

//UpdateMenu Update Menu element object in DB
func (r MenusAPIServer) UpdateMenu(ctx context.Context, menu *pbm.Menu) (*pbm.SimpleResponse, error) {
	modelMenu := ConvertMenuPMtoModel(menu)
	menuResponse, err := r.MenusService.UpdateMenu(modelMenu)
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: menuResponse.Success, Message: menuResponse.Message}, nil
}

//DeleteMenu Delete Menu element Object from DB using ScheduleID
func (r MenusAPIServer) DeleteMenu(ctx context.Context, request *pbm.MenuRequest) (*pbm.SimpleResponse, error) {
	err := r.MenusService.DeleteMenu(request.Id, request.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pbm.SimpleResponse{Success: true, Message: "menu item correctly deleted"}, nil
}

func NewMenusAPIServer(ctx context.Context) *MenusAPIServer {
	return &MenusAPIServer{
		MenusService: services.NewMenuService(ctx),
	}
}

// ConvertRoutePMtoModel Convert Menu PBM to Model object
func ConvertMenuPMtoModel(pbmMenu *pbm.Menu) *models.Menu {
	menu := &models.Menu{
		ID:          pbmMenu.Id,
		LoncheraID:  pbmMenu.LoncheraId,
		Name:        pbmMenu.Name,
		Description: pbmMenu.Description,
		Price:       pbmMenu.Price,
		Currency:    pbmMenu.Currency,
		ImageURL:    pbmMenu.ImageUrl,
	}

	return menu
}

// ConvertMenuModelToPM Convert Menu Object to PBM
func ConvertMenuModelToPM(menu *models.Menu) *pbm.Menu {
	menuModel := &pbm.Menu{
		Id:          menu.ID,
		LoncheraId:  menu.LoncheraID,
		Name:        menu.Name,
		Description: menu.Description,
		Price:       menu.Price,
		Currency:    menu.Currency,
		ImageUrl:    menu.ImageURL,
		CreatedAt:   menu.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   menu.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return menuModel

}
