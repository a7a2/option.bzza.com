define('../../settings/defaultGraphicWindowSettings', ['base'],
    function (WT) {
        WT.getDefaultColorSchemes = function () {
            //additional default color schemes might be created
            var DefaultDark = {
                outerBackground: '#141414',
                innerBackground: '#141414',
                text: '#646464',
                grid: '#232323',
                names: '#7d7d7d',
                line: '#63bf31',
                barUp: '#63bf31',
                barDown: '#dd3b3b',
                bearCandle: '#dd3b3b',
                bullCandle: '#63bf31',
                askLine: '#b2b2b2',
                bidLine: '#b2b2b2',
                currentPriceText: '#333333',
                cross: '#e6e6e6',
                crossLine: '#333333',
                volume: '#63bf31',
                positionVolume: '#ffffff',
                positionLineBuy: '#539f29',
                positionLineSell: '#dd3b3b',
                pendingOrderLineBuy: '#539f29',
                pendingOrderLineSell: '#dd3b3b',
                takeProfitLineBuy: '#539f29',
                takeProfitLineSell: '#dd3b3b',
                stopLossLineBuy: '#539f29',
                stopLossLineSell: '#dd3b3b',
                openDateLine: '#63bf31',
                stopLine: '#ff812d',
                expiryLine: '#ff812d',
                highOptionPositionLine: '#437c24',
                lowOptionPositionLine: '#ab3131',
                topRangeOptionLine: '#437c24',
                bottomRangeOptionLine: '#ab3131',
                topOneTouchOptionLine: '#437c24',
                bottomOneTouchOptionLine: '#ab3131',
                textOfPriceLevel: '#ffffff',
                overlayButtons: '#ffffff'
            };
            var DefaultLight = {
                outerBackground: '#fafafa',
                innerBackground: '#fafafa',
                text: '#969696',
                grid: '#e6e6e6',
                names: '#979695',
                line: '#63bf31',
                barUp: '#63bf31',
                barDown: '#dd3b3b',
                bearCandle: '#dd3b3b',
                bullCandle: '#63bf31',
                askLine: '#3c3c3c',
                bidLine: '#3c3c3c',
                currentPriceText: '#fafafa',
                cross: '#323232',
                crossLine: '#ddddda',
                volume: '#63bf31',
                positionVolume: '#333333',
                positionLineBuy: '#539f29',
                positionLineSell: '#dd3b3b',
                pendingOrderLineBuy: '#539f29',
                pendingOrderLineSell: '#dd3b3b',
                takeProfitLineBuy: '#539f29',
                takeProfitLineSell: '#dd3b3b',
                stopLossLineBuy: '#539f29',
                stopLossLineSell: '#dd3b3b',
                openDateLine: '#63bf31',
                stopLine: '#ff812d',
                expiryLine: '#ff812d',
                highOptionPositionLine: '#437c24',
                lowOptionPositionLine: '#ab3131',
                topRangeOptionLine: '#437c24',
                bottomRangeOptionLine: '#ab3131',
                topOneTouchOptionLine: '#437c24',
                bottomOneTouchOptionLine: '#ab3131',
                textOfPriceLevel: '#ffffff',
                overlayButtons: '#000000'
            };
            var CoffeeDark = {
                outerBackground: '#201e1c',
                innerBackground: '#201e1c',
                text: '#746454',
                grid: '#322c28',
                names: '#7d7d7d',
                line: '#e0924c',
                barUp: '#e0924c',
                barDown: '#dd3b3b',
                bearCandle: '#dd3b3b',
                bullCandle: '#e0924c',
                askLine: '#beb5ad',
                bidLine: '#beb5ad',
                currentPriceText: '#333333',
                cross: '#e6e6e6',
                crossLine: '#333333',
                volume: '#e0924c',
                positionVolume: '#ffffff',
                positionLineBuy: '#539f29',
                positionLineSell: '#dd3b3b',
                pendingOrderLineBuy: '#539f29',
                pendingOrderLineSell: '#dd3b3b',
                takeProfitLineBuy: '#539f29',
                takeProfitLineSell: '#dd3b3b',
                stopLossLineBuy: '#539f29',
                stopLossLineSell: '#dd3b3b',
                openDateLine: '#63bf31',
                stopLine: '#ff812d',
                expiryLine: '#ff812d',
                highOptionPositionLine: '#437c24',
                lowOptionPositionLine: '#ab3131',
                topRangeOptionLine: '#437c24',
                bottomRangeOptionLine: '#ab3131',
                topOneTouchOptionLine: '#437c24',
                bottomOneTouchOptionLine: '#ab3131',
                textOfPriceLevel: '#ffffff',
                overlayButtons: '#ffffff'
            };
            var CoffeeLight = {
                outerBackground: '#fbfaf9',
                innerBackground: '#fbfaf9',
                text: '#a69686',
                grid: '#eae6e2',
                names: '#979695',
                line: '#e0924c',
                barUp: '#e0924c',
                barDown: '#dd3b3b',
                bearCandle: '#dd3b3b',
                bullCandle: '#e0924c',
                askLine: '#463c32',
                bidLine: '#463c32',
                currentPriceText: '#fafafa',
                cross: '#323232',
                crossLine: '#ddddda',
                volume: '#e0924c',
                positionVolume: '#333333',
                positionLineBuy: '#539f29',
                positionLineSell: '#dd3b3b',
                pendingOrderLineBuy: '#539f29',
                pendingOrderLineSell: '#dd3b3b',
                takeProfitLineBuy: '#539f29',
                takeProfitLineSell: '#dd3b3b',
                stopLossLineBuy: '#539f29',
                stopLossLineSell: '#dd3b3b',
                openDateLine: '#63bf31',
                stopLine: '#ff812d',
                expiryLine: '#ff812d',
                highOptionPositionLine: '#437c24',
                lowOptionPositionLine: '#ab3131',
                topRangeOptionLine: '#437c24',
                bottomRangeOptionLine: '#ab3131',
                topOneTouchOptionLine: '#437c24',
                bottomOneTouchOptionLine: '#ab3131',
                textOfPriceLevel: '#ffffff',
                overlayButtons: '#000000'
            };
            var ContrastDark = {
                outerBackground: '#141414',
                innerBackground: '#141414',
                text: '#969696',
                grid: '#414141',
                names: '#7d7d7d',
                line: '#e6e6e6',
                barUp: '#e6e6e6',
                barDown: '#e6e6e6',
                bearCandle: '#e6e6e6',
                bullCandle: '#141414',
                askLine: '#b5b5b5',
                bidLine: '#b5b5b5',
                currentPriceText: '#333333',
                cross: '#e6e6e6',
                crossLine: '#333333',
                volume: '#e6e6e6',
                positionVolume: '#ffffff',
                positionLineBuy: '#539f29',
                positionLineSell: '#dd3b3b',
                pendingOrderLineBuy: '#539f29',
                pendingOrderLineSell: '#dd3b3b',
                takeProfitLineBuy: '#539f29',
                takeProfitLineSell: '#dd3b3b',
                stopLossLineBuy: '#539f29',
                stopLossLineSell: '#dd3b3b',
                openDateLine: '#63bf31',
                stopLine: '#ff812d',
                expiryLine: '#ff812d',
                highOptionPositionLine: '#437c24',
                lowOptionPositionLine: '#ab3131',
                topRangeOptionLine: '#437c24',
                bottomRangeOptionLine: '#ab3131',
                topOneTouchOptionLine: '#437c24',
                bottomOneTouchOptionLine: '#ab3131',
                textOfPriceLevel: '#ffffff',
                overlayButtons: '#ffffff'
            };
            var ContrastLight = {
                outerBackground: '#fafafa',
                innerBackground: '#fafafa',
                text: '#646464',
                grid: '#c8c8c8',
                names: '#979695',
                line: '#141414',
                barUp: '#141414',
                barDown: '#141414',
                bearCandle: '#141414',
                bullCandle: '#fafafa',
                askLine: '#141414',
                bidLine: '#141414',
                currentPriceText: '#fafafa',
                cross: '#323232',
                crossLine: '#ddddda',
                volume: '#141414',
                positionVolume: '#333333',
                positionLineBuy: '#539f29',
                positionLineSell: '#dd3b3b',
                pendingOrderLineBuy: '#539f29',
                pendingOrderLineSell: '#dd3b3b',
                takeProfitLineBuy: '#539f29',
                takeProfitLineSell: '#dd3b3b',
                stopLossLineBuy: '#539f29',
                stopLossLineSell: '#dd3b3b',
                openDateLine: '#63bf31',
                stopLine: '#ff812d',
                expiryLine: '#ff812d',
                highOptionPositionLine: '#437c24',
                lowOptionPositionLine: '#ab3131',
                topRangeOptionLine: '#437c24',
                bottomRangeOptionLine: '#ab3131',
                topOneTouchOptionLine: '#437c24',
                bottomOneTouchOptionLine: '#ab3131',
                textOfPriceLevel: '#ffffff',
                overlayButtons: '#000000'
            };
            var IceDark = {
                outerBackground: '#1b1c21',
                innerBackground: '#1b1c21',
                text: '#545c74',
                grid: '#282932',
                names: '#7d7d7d',
                line: '#4988e4',
                barUp: '#4988e4',
                barDown: '#dd3b3b',
                bearCandle: '#dd3b3b',
                bullCandle: '#4988e4',
                askLine: '#adb1be',
                bidLine: '#adb1be',
                currentPriceText: '#333333',
                cross: '#e6e6e6',
                crossLine: '#333333',
                volume: '#4988e4',
                positionVolume: '#ffffff',
                positionLineBuy: '#539f29',
                positionLineSell: '#dd3b3b',
                pendingOrderLineBuy: '#539f29',
                pendingOrderLineSell: '#dd3b3b',
                takeProfitLineBuy: '#539f29',
                takeProfitLineSell: '#dd3b3b',
                stopLossLineBuy: '#539f29',
                stopLossLineSell: '#dd3b3b',
                openDateLine: '#63bf31',
                stopLine: '#ff812d',
                expiryLine: '#ff812d',
                highOptionPositionLine: '#437c24',
                lowOptionPositionLine: '#ab3131',
                topRangeOptionLine: '#437c24',
                bottomRangeOptionLine: '#ab3131',
                topOneTouchOptionLine: '#437c24',
                bottomOneTouchOptionLine: '#ab3131',
                textOfPriceLevel: '#ffffff',
                overlayButtons: '#ffffff'
            };
            var IceLight = {
                outerBackground: '#f9fafb',
                innerBackground: '#f9fafb',
                text: '#868ea6',
                grid: '#e2e4ea',
                names: '#979695',
                line: '#4988e4',
                barUp: '#4988e4',
                barDown: '#dd3b3b',
                bearCandle: '#dd3b3b',
                bullCandle: '#4988e4',
                askLine: '#323846',
                bidLine: '#323846',
                currentPriceText: '#fafafa',
                cross: '#323232',
                crossLine: '#ddddda',
                volume: '#4988e4',
                positionVolume: '#333333',
                positionLineBuy: '#539f29',
                positionLineSell: '#dd3b3b',
                pendingOrderLineBuy: '#539f29',
                pendingOrderLineSell: '#dd3b3b',
                takeProfitLineBuy: '#539f29',
                takeProfitLineSell: '#dd3b3b',
                stopLossLineBuy: '#539f29',
                stopLossLineSell: '#dd3b3b',
                openDateLine: '#63bf31',
                stopLine: '#ff812d',
                expiryLine: '#ff812d',
                highOptionPositionLine: '#437c24',
                lowOptionPositionLine: '#ab3131',
                topRangeOptionLine: '#437c24',
                bottomRangeOptionLine: '#ab3131',
                topOneTouchOptionLine: '#437c24',
                bottomOneTouchOptionLine: '#ab3131',
                textOfPriceLevel: '#ffffff',
                overlayButtons: '#000000'
            };
            //It's important to add color scheme's name here to make it available.
            return [DefaultDark, DefaultLight, CoffeeDark, CoffeeLight, ContrastDark, ContrastLight, IceDark, IceLight];
        };

        WT.getColorSchemesNames = function () {
            //Respect the sequences that determined above. The "Custom" always must be the last one.
            return ["DefaultDark", "DefaultLight", "CoffeeDark", "CoffeeLight", "ContrastDark", "ContrastLight", "IceDark", "IceLight", "Custom"];
        };

        WT.getLightColorSchemes = function () {
            var def = WT.getDefaultColorSchemes();
            return [def[1], def[3], def[5], def[7]];
        };

        WT.getDarkColorSchemes = function () {
            var def = WT.getDefaultColorSchemes();
            return [def[0], def[2], def[4], def[6]];
        };

        WT.getDefaultWindowsTemplate = function () {
            return {
                /*Possible codes for period:
                 PC_S5: 805,
                 PC_M1: 101,
                 PC_M15: 115,
                 PC_H1: 201,
                 PC_D1: 301
                 */
                period: 201,
                useOffset: true,
                //possible offsets: 0, 0.05, 0.15, 0.3. This option is useless if 'useOffset' is false.
                offset: 0.3,
                //possible scales: 1, 2, 4, 8, 16, 32
                scale: 8,
                autoScroll: true,
                autoTimeframe: true,
                showGrid: true,
                showNames: true,

                /*Possible codes for barStyle:
                 gBar: 0,
                 jCandle: 1,
                 line: 2
                 */
                barStyle: 1,
                showCurrentPrice: true,
                showAskLine: false,
                showOptionSettingsLine: true,
                //set the number of default color scheme, that should be applied at startup.
                //Numbers start from 0, i.e
                //0 - first,
                //1 - second etc.
                colorScheme: WT.getDefaultColorSchemes()[0]
            };
        };

        WT.getFirstStartCharts = function () {
            //charts that should be shown when user opens WebTerminal first time.
            return ["EURUSD", "GBPUSD", "USDCHF", "USDJPY"];
        };

        WT.getFirstStartActiveChart = function () {
            //chart, that should be active when user opens WebTerminal first time.
            return "GBPUSD";
        };

        WT.showTradePanel = function () {
            //should trade panel be active when user opens charts.
            return false;
        };

        return WT;
    });