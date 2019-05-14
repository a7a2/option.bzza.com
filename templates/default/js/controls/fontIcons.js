    // Объект обеспечивающий более удобный и централизированный доступ к
    // иконкам. Соответственно после обновления шрифтов все новые ионки нужно
    // добавить сюда, или же обновить значения кодов символов, если они изменились в шрифте.
;

(function(window) {
    var func = String.fromCharCode;

    window.FontIcons = {
        // Icons for Trade lines
        SellLine:           func('0xe907'),
        BuyLine:            func('0xe908'),
        BuyLimitLine:       func('0xe904'),
        BuyStopLine:        func('0xe905'),
        SellLimitLine:      func('0xe90b'),
        SellStopLine:       func('0xe90c'),
        HighOptionLine:     func('0xe908'),
        LowOptionLine:      func('0xe907'),
        TopOneTouchLine:    func('0xe90d'),
        BottomOneTouchLine: func('0xe902'),
        InsideRangeLine:    func('0xe928'),
        OutsideRangeLine:   func('0xe927'),

        Eye: func('0xe914'),

        // Currency
        Usd: func('0xe90e'),
        Eur: func('0xe90f'),
        Gbp: func('0xe911'),
        Rub: func('0xe912'),
        Jpy: func('0xe913'),
        Chf: func('0xe910'),

        // m_ - префикс означающий, что эта иконка предназначена для мобильного приложения, но не обязательно))))
        // Trade buttons
        m_NewOrder:     func('0xe94a'),
        m_BuyOption:    func('0xe948'),
        m_PendingOrder: func('0xe949'),

        // Swiper page icons
        m_SettingsPage: func('0xe947'),
        m_TradePage:    func('0xe946'),
        m_ReportPage:   func('0xe945'),
        m_QuotesPage:   func('0xe944'),

        // Chart styles
        m_BullCandle: func('0xe934'),
        m_BearCandle: func('0xe935'),
        m_LineCandle: func('0xe936'),

        // etc
        m_Hyphen:          func('0xe940'),
        m_GoToBack:        func('0xe93f'),
        m_PasswordEye:     func('0xe93e'),
        m_EmptyQuotes:     func('0xe93d'),
        m_EmptyCheckBox:   func('0xe93c'),
        m_ShadedArrowDown: func('0xe93b'),
        m_ArrowDown:       func('0xe93a'),
        m_CrossedEye:      func('0xe939'),
        m_CircleCheck:     func('0xe938'),
        m_CheckBox:        func('0xe937'),
        m_Star:            func('0xe94d'),
        m_TradePanel:      func('0xe94c'),
        m_Indicators:      func('0xe94b'),
        m_ExitCross:       func('0xe931'),
        m_Params:          func('0xe94f'),
        m_RadioCheck:      func('0xe94e'),
        m_Radio:           func('0xe952'),
        m_Aim:             func('0xe951'),
        m_OptionTypes:     func('0xe950'),

        m_minus:           func('0xe953'),
        m_plus:            func('0xe954'),
        m_tableTradePlus:  func('0xe956'),
        m_Download:        func('0xe955'),


        getIconByChartCode: function(chartCode) {
            var result = '';
            try {
                result = String.fromCharCode(chartCode);
            } catch (ex) {
                console.log('Invalid chart code');
                result = ''
            }
            return result;
        },

        getIconByDepositCurrency: function (WT) {
            var traderGroup = WT.TraderData.getTraderGroup();
            if (!traderGroup) {
                return;
            }

            var depositCurrency = traderGroup.depositCurrency.toLowerCase();
            switch (depositCurrency) {
                case 'usd':
                    return WT.FontIcons.Usd;
                    break;
                case 'eur':
                    return WT.FontIcons.Eur;
                    break;
                case 'gbp':
                    return WT.FontIcons.Gbp;
                    break;
                case 'rub':
                    return WT.FontIcons.Rub;
                    break;
                case 'jpy':
                    return WT.FontIcons.Jpy;
                    break;
                case 'chf':
                    return WT.FontIcons.Chf;
                    break;
            }
        }
    };
})(window);