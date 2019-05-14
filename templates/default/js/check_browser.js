if (!("WebSocket" in window))
{
    document.location.href = "unsupportedBrowser.html";
}
var browser, browserSplit, UA,
    OperaB, OperaV, Opera,
    FirefoxB, Firefox,
    ChromeB, ChromeiOSB, Chrome, ChromeiOS,
    SafariB, SafariV, Safari,
    IEB, IE, IEB2, IE2,
    IsiPad, IsiPhone, IsMac;
UA = window.navigator.userAgent;// содержит переданный браузером юзерагент
// шаблоны для распарсивания юзерагента

OperaB = /Opera[ \/]+\w+\.\w+/i;
OperaV = /Version[ \/]+\w+\.\w+/i;
FirefoxB = /Firefox\/\w+\.\w+/i;
ChromeB = /Chrome\/\w+\.\w+/i;
ChromeiOSB = /CriOS\/\w+\.\w+/i;
SafariB = /Version\/\w+\.\w+/i;
IEB = /MSIE *\d+\.\w+/i;
IEB2 = /Trident\/\d/i;
SafariV = /Safari\/\w+\.\w+/i;
browser = []; //массив с данными о браузере
browserSplit = /[ \/\.]/i; //шаблон для разбивки данных о браузере из строки
OperaV = UA.match(OperaV);
Firefox = UA.match(FirefoxB);
Chrome = UA.match(ChromeB);
ChromeiOS = UA.match(ChromeiOSB);
Safari = UA.match(SafariB);
SafariV = UA.match(SafariV);
IE = UA.match(IEB);
IE2 = UA.match(IEB2);
Opera = UA.match(OperaB);
IsiPad = UA.indexOf("iPad") > -1;
IsiPhone = UA.indexOf("iPhone") > -1;
IsMac = UA.indexOf("Mac") > -1;

if (!ChromeiOS == "")
    browser[0] = ChromeiOS[0];
else if ((!Opera == "") && (!OperaV == ""))
    browser[0] = OperaV[0].replace(/Version/, "Opera");
else if (!Opera == "")
    browser[0] = Opera[0];
else if (!IE2 == "")
    browser[0] = IE2[0];
else if (!IE == "")
    browser[0] = IE[0];
else if (!Firefox == "")
    browser[0] = Firefox[0];
else if (!Chrome == "")
    browser[0] = Chrome[0];
else if ((!Safari == "") && (!SafariV == ""))
    browser[0] = Safari[0].replace("Version", "Safari");

//------------ Разбивка версии -----------
var data = []; // [0] - имя браузера, [1] - целая часть версии
if (browser[0] != null)
    data = browser[0].split(browserSplit);
var checked = false,
    version = parseInt(data[1], 10);
browser = data[0];

var isMobile =
{
    Android: UA.match(/Android/i) || false,
    BlackBerry: UA.match(/BlackBerry/i) || false,
    iOS: UA.match(/iPhone|iPad|iPod/i) || false,
    Opera: UA.match(/Opera Mini/i) || false,
    Windows: (UA.match(/Trident/i) && UA.match(/ARM/i)) || false
};
isMobile.any = isMobile.Android || isMobile.BlackBerry || isMobile.iOS || isMobile.Opera || isMobile.Windows || false;
switch (browser) {
    case "CriOS": // Chrome под iOS
        checked = version >= 23;
        break;
    case "Chrome":
        checked = version >= 23;
        break;
    case "Trident":
        checked = version >= 6;
        break;
    case "MSIE":
        checked = version >= 10;
        break;
    case "Safari":
        checked = IsMac && version >= 8;
        break;
    case "Firefox":
        checked = version >= 10;
        break;
    case "Opera":
        checked = version >= 15;
        break;
    default:
        checked = false;
        break;
}

// 0 - настольная версия,
// 1 - универсальная настройка (веб-терминал, планшетная версия, мобильная версия),
// 2 - мобильное приложение.
var webTerminalMode = null,
    wtmConfigProperty = window.UtipConfig.TerminalSettings.WebTerminalMode;
if (wtmConfigProperty !== null && wtmConfigProperty !== undefined
    && typeof wtmConfigProperty !== 'string' && !isNaN(wtmConfigProperty)) {
    webTerminalMode = wtmConfigProperty;
} else {
    ///
    // по стандарту выставляем режим "универсальная настройка (веб-терминал, планшетная версия, мобильная версия)",
    // т. к. он дает возможность корректно зайти на любую платформу и выбирать сервера с конфига или удаленки,
    // исключая веб-терминал - он может выбирать только сервера из конфига. не подключает скрипт cordova.
    ///
    webTerminalMode = 1;
}

if (webTerminalMode === 0 && isMobile.any) {
    checked = true;
    if (UA.match(/Windows Phone 10.0/i) != null) {
        window.location.href = "unsupportedOther.html";
    } else if ((isMobile.Android && UA.toLowerCase().indexOf("windows phone") == -1) || (browser == "Chrome" && UA.match(/Edge/i) == null)) {
        window.location.href = "unsupportedAndroid.html";
    } else if ((isMobile.iOS && UA.toLowerCase().indexOf("windows phone") == -1) || (browser == "CriOS" || browser == "Safari")) {
        window.location.href = "unsupportedIOS.html";
    } else if (isMobile.Windows || isMobile.Opera) {
        window.location.href = "unsupportedOther.html";
    }
}

var find = function (needle) {
    return UA.toLowerCase().indexOf(needle) !== -1;
};
var WT_Device = {
    ios:                function () { return WT_Device.iphone() || WT_Device.ipod() || WT_Device.ipad(); },
    iphone:             function () { return !WT_Device.windows() && find('iphone'); },
    ipod:               function () { return find('ipod'); },
    ipad:               function () { return find('ipad'); },
    android:            function () { return !WT_Device.windows() && find('android'); },
    androidPhone:       function () { return WT_Device.android() && find('mobile'); },
    androidTablet:      function () { return WT_Device.android() && !find('mobile'); },
    blackberry:         function () { return find('blackberry') || find('bb10') || find('rim'); },
    blackberryPhone:    function () { return WT_Device.blackberry() && !find('tablet'); },
    blackberryTablet:   function () { return WT_Device.blackberry() && !find('tablet'); },
    windows:            function () { return find('windows'); },
    windowsPhone:       function () { return (WT_Device.windows() && find('phone')) || (find('trident') && find('arm')); },
    windowsTablet:      function () { return WT_Device.windows() && (find('touch') && !WT_Device.windowsPhone()); },
    fxos:               function () { return (find('(mobile;') || find('(tablet;')) && find('; rv:'); },
    fxosPhone:          function () { return WT_Device.fxos() && find('mobile'); },
    fxosTablet:         function () { return WT_Device.fxos() && find('tablet'); },
    meego:              function () { return find('meego'); },
    cordova:            function () { return window.cordova && location.protocol === 'file:'; },
    nodeWebkit:         function () { return typeof window.process === 'object'; },
    mobile:             function () {
                            return WT_Device.androidPhone() || WT_Device.iphone() || WT_Device.ipod() || WT_Device.windowsPhone()
                                || WT_Device.blackberryPhone() || WT_Device.fxosPhone() || WT_Device.meego();
                        },
    tablet:             function () {
                            return WT_Device.ipad() || WT_Device.androidTablet() || WT_Device.blackberryTablet()
                                || WT_Device.windowsTablet() || WT_Device.fxosTablet();
                        },
    desktop:            function () { return !WT_Device.tablet() && !WT_Device.mobile(); }
};

function IsTabletApplication() {
    return !WT_Device.mobile() && !WT_Device.desktop() && WT_Device.tablet();
}

//определяем запуск на мобильнике или большом экране и подтягиваем нужную CSS-ку.
var cellPhone = false;
var link = document.createElement('link');
link.setAttribute('type','text/css');
link.setAttribute('rel','stylesheet');
if (WT_Device.mobile() && webTerminalMode !== 0) {
    cellPhone = true;
    link.setAttribute('href',' css/wtMobile.css?v=1.11.4.1');
} else {
    var settings = JSON.parse(localStorage.getItem('UserSettings'));
    if (!settings) {
        link.setAttribute('href', ' css/wtDarkStyle.css?v=1.11.4.1');
    } else {
        var colorScheme = settings.Terminal.ColorScheme;

        switch (colorScheme) {
            case 'black_color_scheme':
                link.setAttribute('href', ' css/wtDarkStyle.css?v=1.11.4.1');
                break;
            case 'beige_color_scheme':
                link.setAttribute('href', ' css/wtStyle.css?v=1.11.4.1');
                break;
            default:
                link.setAttribute('href', ' css/wtDarkStyle.css?v=1.11.4.1');
                break;
        }
    }
}
document.getElementsByTagName('head')[0].appendChild(link);

var serversMode = null,
    configProperty = window.UtipConfig.TerminalSettings.ServersMode;
if (configProperty !== null && configProperty !== undefined && (cellPhone === true || WT_Device.tablet())) {
    serversMode = configProperty;
} else {
    serversMode = webTerminalMode;
}

if (isMobile.iOS && webTerminalMode !== 0) {
    //всегда разрешаем запускаться приложению на iPhone, целевая версия iOS для приложения всё равно стоит не ниже 8.0
    checked = true;
    //подключаем поддержку поворота (iOS fix).
    window.shouldRotateToOrientation = function(degrees) {
        return true;
    }
}
if (webTerminalMode == 2) {
    var cordovaScript = document.createElement('script');
    cordovaScript.setAttribute('src', 'cordova.js');
    cordovaScript.setAttribute('type', 'application/javascript');
    document.getElementsByTagName('head')[0].appendChild(cordovaScript);
}


var Browser = {
    IsMobile: cellPhone,
    Browser: browser,
    Version: version,
    IsIPhone: WT_Device.iphone() && !WT_Device.windowsPhone(),
    IsTablet: WT_Device.tablet()
};

if (!checked) {
    document.location.href = "unsupportedBrowser.html";
}