//должен быть доступен отовсюду. В том числе и до загрузки модуля терминала.
var preloader;
var isMobile = Browser.IsMobile;
var preloaderIsShowed = false;
var preloaderType;

function createPreloader () {
    var logoImage = new Image();
    logoImage.src = "res/img/intro_logo.png";
    var mobileLogoImage = new Image();
    mobileLogoImage.src = "res/img/mobileImg/intro_logo.png";
    var circleImage = new Image();
    circleImage.src = "res/img/intro_spinner.png";


    preloader = document.createElement("div");
    preloader.id = "preloader";
    preloader.classList.add('preloader');


    var imageBlock = document.createElement("div");
    imageBlock.id = "imageBlock";

    var preloaderLogo = document.createElement("img");
    preloaderLogo.id = "preloaderLogo";
    preloaderLogo.src = !isMobile ? "res/img/intro_logo.png" : "res/img/mobileImg/intro_logo.png";
    imageBlock.appendChild(preloaderLogo);

    var preloaderCircle = document.createElement("img");
    preloaderCircle.id = "preloaderCircle";
    preloaderCircle.src = "res/img/intro_spinner.png";
    imageBlock.appendChild(preloaderCircle);

    preloader.appendChild(imageBlock);

    tuneStyles(preloader, imageBlock, preloaderLogo, preloaderCircle);

    document.body.appendChild(preloader);
}
//настраиваем стили прямо здесь, чтобы не быть зависимыми от недогруженных файлов.
function tuneStyles (preloader, imageBlock, preloaderLogo, preloaderCircle) {
    var styleSheetElement = document.createElement("style"), customStyleSheet;
    document.head.appendChild(styleSheetElement);
    customStyleSheet = document.styleSheets[document.styleSheets.length-1];

    preloader.style.position = "fixed";
    preloader.style.zIndex = "13";
    preloader.style.width = "200px";
    preloader.style.height = "200px";
    preloader.style.display = 'none';
    preloader.style.opacity = '0';

    imageBlock.style.position = "relative";
    imageBlock.style.display = "inline-block";
    imageBlock.style.width = "200px";
    imageBlock.style.height = "200px";

    preloaderCircle.style.display = "none";
    preloaderCircle.style.width = "200px";
    preloaderCircle.style.height = "200px";
    preloaderCircle.style.position = "absolute";
    preloaderCircle.classList.add("preloaderCircle");

    preloaderLogo.style.display = "none";
    preloaderLogo.style.position = "absolute";
    preloaderLogo.style.top = "36px";
    preloaderLogo.style.left = !isMobile ? "30px" : "51px";
}

var freezed = false,
    fn = function (e) { e.preventDefault(); }.bind(this);

function freezeHTML () {
    if (!freezed) {
        document.documentElement.addEventListener('touchmove', fn);
        document.documentElement.style.pointerEvents = 'none';
        document.documentElement.scrollTop = 0;
        freezed = true;
    }
}
function unfreezeHTML () {
    if (freezed) {
        document.documentElement.removeEventListener('touchmove', fn);
        document.documentElement.style.pointerEvents = '';
        freezed = false;
    }
}

function showPreloader () {
    if (window.AuthModule && window.AuthModule.showed === true) return;

    preloader.style.display = "block";
    preloader.style.opacity = '1';
    preloader.querySelector("#preloaderLogo").style.display = "block";
    preloader.querySelector("#preloaderCircle").style.display = "block";

    preloaderType = 'logo_circle';
    preloaderIsShowed = true;

    if (WT_Device.iphone()) {
        freezeHTML();
    }
}
function showPreloaderCircle () {
    if (window.AuthModule && window.AuthModule.showed === true) return;

    preloader.style.display = "block";
	preloader.style.opacity = "1";
    preloaderType = 'circle';
    preloaderIsShowed = true;

    preloader.querySelector("#preloaderCircle").style.display = "block";

    if (WT_Device.iphone()) {
        freezeHTML();
    }
}
function hidePreloader () {
    preloader.style.opacity = '0';
    setTimeout(function () {
        preloader.style.display = 'none';
        preloader.querySelector("#preloaderCircle").style.display = "none";
        preloader.querySelector("#preloaderLogo").style.display = "none";

        preloaderType = '';
        preloaderIsShowed = false;

        if (WT_Device.iphone()) {
            unfreezeHTML();
        }
    }, 250);
}
createPreloader();
