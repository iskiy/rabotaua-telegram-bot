package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	schedulesText             string
	viewVacanciesButton       tgbotapi.KeyboardButton
	addParametersButton       tgbotapi.KeyboardButton
	addSubscriptionsButton    tgbotapi.KeyboardButton
	deleteParametersButton    tgbotapi.KeyboardButton
	deleteSubscriptionsButton tgbotapi.KeyboardButton
	cancelButton              tgbotapi.KeyboardButton
	whateverButton            tgbotapi.KeyboardButton

	mainMenuKeyboard   tgbotapi.ReplyKeyboardMarkup
	cancelKeyboard     tgbotapi.ReplyKeyboardMarkup
	parametersKeyboard tgbotapi.ReplyKeyboardMarkup
)

type Config struct {
	Buttons Buttons `yaml:"buttons_text"`
	Msg     Msg     `yaml:"msg"`
}

type Buttons struct {
	ViewVacancies      string `yaml:"view_vacancies"`
	AddParameters      string `yaml:"add_parameters"`
	AddSubscriptions   string `yaml:"add_subscriptions"`
	DeleteParameters   string `yaml:"delete_parameters"`
	DeleteSubsriptions string `yaml:"delete_subsriptions"`
	Cancel             string `yaml:"cancel"`
	Whatever           string `yaml:"whatever"`
}

type Msg struct {
	Start                  string `yaml:"start"`
	UndefinedCommand       string `yaml:"undefined_command"`
	MainMenu               string `yaml:"main_menu"`
	EnterPostname          string `yaml:"enter_postname"`
	Undefined              string `yaml:"undefined"`
	WrongSubNum            string `yaml:"wrong_sub_num"`
	WrongParamsNum         string `yaml:"wrong_params_num"`
	AddParamsSuccess       string `yaml:"add_params_success"`
	EmptyParameters        string `yaml:"empty_parameters"`
	EnterParamNum          string `yaml:"enter_param_num"`
	UndefinedView          string `yaml:"undefined_view"`
	LeftBorder             string `yaml:"left_border"`
	RightBorder            string `yaml:"right_border"`
	DeleteParameterSuccess string `yaml:"delete_parameter_success"`
	AlreadySub             string `yaml:"already_sub"`
	AddSubSuccess          string `yaml:"add_sub_success"`
	EmptySubs              string `yaml:"empty_subs"`
	EnterSubsNum           string `yaml:"enter_subs_num"`
	DeleteSubSuccess       string `yaml:"delete_sub_success"`
	FoundSub               string `yaml:"found_sub"`
	CantFindVacancies      string `yaml:"cant_find_vacancies"`
	WrongEmploymentTypeNum string `yaml:"wrong_employment_type_num"`
	EnterEmploymentTypeNum string `yaml:"enter_employment_type_num"`
	UndefinedCity          string `yaml:"undefined_city"`
	EnterCity              string `yaml:"enter_city"`
}

var cfg Config

func Init() error {
	err := initConfig()
	if err != nil {
		return err
	}
	initButtons()
	initKeyboards()
	return nil
}

func initConfig() error {
	data, err := os.ReadFile("configs/config.yml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}
	return nil
}

func initButtons() {
	viewVacanciesButton = tgbotapi.NewKeyboardButton(cfg.Buttons.ViewVacancies)
	addParametersButton = tgbotapi.NewKeyboardButton(cfg.Buttons.AddParameters)
	addSubscriptionsButton = tgbotapi.NewKeyboardButton(cfg.Buttons.AddSubscriptions)
	deleteParametersButton = tgbotapi.NewKeyboardButton(cfg.Buttons.DeleteParameters)
	deleteSubscriptionsButton = tgbotapi.NewKeyboardButton(cfg.Buttons.DeleteSubsriptions)
	cancelButton = tgbotapi.NewKeyboardButton(cfg.Buttons.Cancel)
	whateverButton = tgbotapi.NewKeyboardButton(cfg.Buttons.Whatever)
}

func initKeyboards() {
	mainMenuKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			viewVacanciesButton, addParametersButton,
		),
		tgbotapi.NewKeyboardButtonRow(
			addSubscriptionsButton, deleteSubscriptionsButton, deleteParametersButton,
		),
	)

	cancelKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(cancelButton),
	)

	parametersKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(whateverButton, cancelButton),
	)
}
