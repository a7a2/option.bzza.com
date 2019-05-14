window.UtipConfig = {
    Servers: [
		{
            Name: 'Binary Demo',
            Url: 'https://Binary1.bzza.com', //перед url обязательно указывается протокол работы сервера: "http://" или "https://"
            Port: '443',
            Type: 'Demo'
        }
    ],
    Groups: [
		{
            Id: 2,
            Name: 'Demo'
        }
    ],
    Language: {
        en: 'English',
        ru: 'Русский',
        fr: 'Français',
        ar: 'العربية',
        pl: 'Polski',
        gr: 'ქართული',
        zh: '中文',
        es: 'Español',
        fa: 'فارسی',
        de: 'Deutsch',
		ko: '한국어'
    },
    BrokersSettings: {
        TerminalName: "WebTrader",
        CompanyName : "bzza",
        CompanySite: "www.bzza.com",
        CompanyMail: "info@a7a2.com",
        Copyrights: "&copy; 2006-2018, bzza.com Technologies Ltd.",
        CopyrightsSite: "www.bzza.com",
        LiveAccountCaption: {
            en: "Open Live Account",
            ru: "Открыть реальный счет",
            fr: "Open Live Account",
            ar: "Open Live Account",
            pl: "Otwórz prawdziwe konto",
            gr: "Open Live Account",
            zh: 'Open Live Account',
            es: 'Open Live Account',
            fa: 'Open Live Account',
            de: 'Open Live Account',
			ko: '실거래 계좌 개설'
        },
        OpenDemoAccountCaption: {
            en: "Open a demo account",
            ru: "Открыть учебный счет",
            fr: "Ouvrir un compte de formation",
            ar: "فتح حساب تجريبي",
            pl: "otwórz rachunek demo",
            gr: "სასწავლო ანგარიშის გახსნა",
            zh: "开通模拟账号",
            es: "Abrir demo cuenta",
            fa: 'حساب دمو افتتاح کنید',
            de: 'Demokonto eröffnen',
			ko: '데모거래 계좌 개설'
        },
        DepositMoneyCaption: { //пополнение счёта
            en: "Deposit money",
            ru: "Пополнить счет",
            fr: "Make a Deposit",
            ar: "Make a Deposit",
            pl: "Make a Deposit",
            gr: "Make a Deposit",
            zh: 'Make a Deposit',
            es: 'Make a Deposit',
            fa: 'Make a Deposit',
            de: 'Einzahlen',
			ko: '입금'
        },
        DepositMoneyLink: "www.bzza.com", //ссылка на пополнение счёта
        //Если такая настройка отсутствует то алгоритм такой:
        //  Если есть брокерская настройка "CompanySite" то при нажатии на кнопку открытия счета подставится эта ссылка
        //  Если нет настройки "CompanySite" то подставляется ссылка по умолчанию "www.bzza.com"
        LiveAccountLink: "www.bzza.com",
        // Если есть брокерская настройка "OpenDemoLink", то при открытии демо-счета будет открыта эта ссылка в новой вкладке
        // ссылка должна начинаться с "http://"
        // Если нет настройки "openDemoLink", то в терминале откроется диалоговое окно с открытием демо-счета (по умолчанию)
        OpenDemoLink: "",
        // Если брокерская настройка "isVisibleOpenDemo" имеет значение "true" или она будет отсутсвовать, то будет открыта возможность открытия демо счета
        // Если брокерская настройка "isVisibleOpenDemo" имеет значение "false" то такая возможность будет закрыта
        isVisibleOpenDemo: true,
        WebOfficeLink: "http://gmtoffice.lh",
        WebOfficeAPILink: "http://weboffice.lh",
        WebOfficeAPIKey: "Fiugkjyu76fhjt7hbk",
        ServerForDemo: "Arsenio", //Должно совпадать с именем сервера в личном кабинете
        DemoGroupId: 2, //id торговой группы из торговой платформы (firebird)
        AdditionalRegFields: [                  //дополнительные поля регистрации
          /*{
                name: "Patronymic",     //Имя поля (помещается в название по-умолчанию)
                key: "cWCTAccountListFormPatr",                //Ключ строкового ресурса (для локализации названия поля)
                parameter: "patronymic", //Название отсылаемого параметра API в личный кабинет
                required: false          //Если true, то поле обязательно для заполнения
            },
            {
                name: "Country",
                key: "cWCTAccountListFormCountry",
                parameter: "country",
                required: false
            },
            {
                name: "Area",
                key: "cWCTAccountListFormRegion",
                parameter: "area",
                required: false
            },
            {
                name: "City",
                key: "cWCTAccountListFormCity",
                parameter: "city",
                required: false
            },
            {
                name: "Address",
                key: "cWCTAccountListFormAddress",
                parameter: "address",
                required: false
            },
            {
                name: "Postcode",
                key: "cWCTAccountListFormIndex",
                parameter: "postcode",
                required: false
            },
            {
                name: "Passport",
                key: "cMobilePassport",
                parameter: "passport",
                required: false
            }*/
        ],
        ProfileMenuItems: [
            {
                Text: "Make a deposit",
                Translations: {
                    en: "Make a deposit",
                    ru: "Внести средства",
                    fr: "Make a deposit",
                    ar: "Make a deposit",
                    pl: "Make a deposit",
                    gr: "Make a deposit",
                    zh: 'Make a deposit',
                    es: 'Make a deposit',
                    fa: 'Make a deposit',
                    de: 'Make a Deposit'
                },
                Id: "payment",
                WebOfficePage: "payment", //по-умолчанию пробуем отобразить раздел веб-офиса в терминале, если поле отсутствует, открываем ссылку в отдельном окне
                Link: "http://www.bzza.com" //указание протокола обязательно
            },
            {
                Text: "My profile",
                Translations: {
                    en: "My profile",
                    ru: "Мой профиль",
                    fr: "My profile",
                    ar: "My profile",
                    pl: "My profile",
                    gr: "My profile",
                    zh: 'My profile',
                    es: 'My profile',
                    fa: 'My profile',
                    de: 'My profile'
                },
                Id: "profile",
                WebOfficePage: "profile",
                Link: "http://www.bzza.com" //указание протокола обязательно
            }
        ],
        Default_GroupId: 2,
        Default_Leverage: 100,
        Default_FirstDeposit: 5000
    },
    TradeSettings: {
        //когда true, то в качестве MaxDeviation используется StopLevel инструмента. В противном случае MaxDeviation = 1000
        UseStopLevel: false,
        DefaultDeposit: 5000
    },
    TerminalSettings: {
        IntegratedWebOfficeEnabled: false,
        ChartFieldHeightPercentage: 61.2,
        SwiperQuotesWidth: 260,
        ServersMode: 0, // режим загрузки списка серверов. 0 - из конфига, 1 - с удаленного сервера
        WebTerminalMode: 1 // 0 - настольная версия, 1 - универсальная версия, 2 - мобильное приложение.
    }
};