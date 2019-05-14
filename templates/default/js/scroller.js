var customScroller = function (containerNode, isMenu) {
    var self = this;
    this.scrollContainer = containerNode;
    this.scrollContentWrapper = containerNode.querySelector('.scroller-content-wrapper');
    this.scrollContent = containerNode.querySelector('.scroller-content');
    this.contentPosition = 0;
    this.scrollerBeingDragged = false;
    this.scroller = null;
    this.topPosition = null;
    this.scrollerHeight = null;

    this.contentLeftPosition = 0;
    this.hScrollerBeingDragged = false;
    this.hScroller = null;
    this.leftPosition = null;
    this.hScrollerWidth = null;
    this.isMenu = !!isMenu;

    this.calculateScrollerHeight = function () {
        var visibleRatio = self.scrollContainer.offsetHeight / self.scrollContent.scrollHeight;
        return visibleRatio * self.scrollContainer.offsetHeight;
    };

    this.calculateScrollerWidth = function () {
        var visibleRatio = self.scrollContainer.offsetWidth / self.scrollContent.scrollWidth;
        return visibleRatio * self.scrollContainer.offsetWidth;
    };

    this.moveScroller = function (evt) {
        self.scrollContentWrapper.scrollTop += Math.sign(evt.deltaY) * 50;
        var scrollPercentage = self.scrollContentWrapper.scrollTop / self.scrollContentWrapper.scrollHeight;
        self.topPosition = scrollPercentage * (self.scrollContainer.offsetHeight - 5); // 5px arbitrary offset so scroll bar doesn't move too far beyond content wrapper bounding box
        self.scroller.style.top = self.topPosition + 'px';

        self.scrollContentWrapper.scrollLeft += Math.sign(evt.deltaX) * 50;
        scrollPercentage = self.scrollContentWrapper.scrollLeft / self.scrollContentWrapper.scrollWidth;
        self.leftPosition = scrollPercentage * (self.scrollContainer.offsetWidth - 5); // 5px arbitrary offset so scroll bar doesn't move too far beyond content wrapper bounding box
        self.hScroller.style.left = self.leftPosition + 'px';
    };

    this.startDrag = function (evt) {
        self.normalizedPosition = evt.pageY;
        self.contentPosition = self.scrollContentWrapper.scrollTop;
        self.scrollerBeingDragged = true;
        window.scrollbar = {
            isMenu: self.isMenu,
            moving: true,
            object: self
        };
    };

    this.startHDrag = function (evt) {
        self.normalizedHPosition = evt.pageX;
        self.contentPosition = self.scrollContentWrapper.scrollLeft;
        self.hScrollerBeingDragged = true;
        window.scrollbar = {
            isMenu: self.isMenu,
            moving: true,
            object: self
        };
    };

    this.touchstart = function(evt) {
        // Special for tablets
        self.startTouchY = evt.targetTouches[0].screenY;
        self.startTouchX = evt.targetTouches[0].screenX; 

        window.scrollbar = {
            isMenu: self.isMenu,
            moving: true,
            object: self
        };
    };

    this.stopDrag = function(evt) {
        self.scrollerBeingDragged = false;
        self.hScrollerBeingDragged = false;
    };

    this.touchMoveScroller = function(evt) {
        evt.preventDefault();
        var deltaY = evt.targetTouches[0].screenY - self.startTouchY,
            deltaX = evt.targetTouches[0].screenX - self.startTouchX;

        // Vertical scrolling ----------------------------------------------
        if (deltaY > self.startTouchY) {
            self.scrollContentWrapper.scrollTop += Math.sign(deltaY) * 25;
        } else if (deltaY < self.startTouchY) {
            self.scrollContentWrapper.scrollTop -= Math.sign(deltaY) * 25;
        }

        var scrollPercentage = self.scrollContentWrapper.scrollTop / self.scrollContentWrapper.scrollHeight;
        self.topPosition = scrollPercentage * (self.scrollContainer.offsetHeight - 5); // 5px arbitrary offset so scroll bar doesn't move too far beyond content wrapper bounding box
        self.scroller.style.top = self.topPosition + 'px';

        // Horizontal scrolling ----------------------------------------------
        if (deltaX > self.startTouchX) {
            self.scrollContentWrapper.scrollLeft += Math.sign(deltaX) * 25;
        } else if (deltaX < self.startTouchX) {
            self.scrollContentWrapper.scrollLeft -= Math.sign(deltaX) * 25;
        }

        scrollPercentage = self.scrollContentWrapper.scrollLeft / self.scrollContentWrapper.scrollWidth;
        self.leftPosition = scrollPercentage * (self.scrollContainer.offsetWidth - 5); // 5px arbitrary offset so scroll bar doesn't move too far beyond content wrapper bounding box
        self.hScroller.style.left = self.leftPosition + 'px';

        self.startTouchY = evt.targetTouches[0].screenY;
        self.startTouchX = evt.targetTouches[0].screenX;
    };


    this.scrollBarScroll = function (evt) {
        var mouseDifferential, scrollEquivalent, scrollPercentage;
        if (self.scrollerBeingDragged === true) {
            mouseDifferential = evt.pageY - self.normalizedPosition;
            scrollEquivalent = mouseDifferential * (self.scrollContentWrapper.scrollHeight / self.scrollContainer.offsetHeight);
            self.scrollContentWrapper.scrollTop = self.contentPosition + scrollEquivalent;
            //move scroller
            scrollPercentage = self.scrollContentWrapper.scrollTop / self.scrollContentWrapper.scrollHeight;
            self.topPosition = scrollPercentage * (self.scrollContainer.offsetHeight - 5); // 5px arbitrary offset so scroll bar doesn't move too far beyond content wrapper bounding box
            self.scroller.style.top = self.topPosition + 'px';
        } else if (self.hScrollerBeingDragged === true) {
            mouseDifferential = evt.pageX - self.normalizedHPosition;
            scrollEquivalent = mouseDifferential * (self.scrollContentWrapper.scrollWidth / self.scrollContainer.offsetWidth);
            self.scrollContentWrapper.scrollLeft = self.contentLeftPosition + scrollEquivalent;
            //move scroller
            scrollPercentage = self.scrollContentWrapper.scrollLeft / self.scrollContentWrapper.scrollWidth;
            self.leftPosition = scrollPercentage * (self.scrollContainer.offsetWidth - 5); // 5px arbitrary offset so scroll bar doesn't move too far beyond content wrapper bounding box
            self.hScroller.style.left = self.leftPosition + 'px';
        }
    };

    this.createScroller = function () {
        self.scroller = document.createElement("div");
        self.scroller.className = 'scroller';
        self.scroller.style.height = '0px';
        self.scrollContainer.appendChild(self.scroller);
        self.scrollContainer.className += ' showScroll';

        //horizontal

        self.hScroller = document.createElement("div");
        self.hScroller.className = 'hScroller';
        self.hScroller.style.width = '0px';
        self.scrollContainer.appendChild(self.hScroller);
        self.scrollContainer.className += ' showHScroll';

        // attach related draggable listeners
        self.scroller.addEventListener('mousedown', self.startDrag);
        self.hScroller.addEventListener('mousedown', self.startHDrag);

        if (self.scrollerHeight / self.scrollContainer.offsetHeight < 1
            || self.hScrollerWidth / self.scrollContainer.offsetWidth < 1) {
            // *If there is a need to have scroll bar based on content size
            self.update();
        }

    };

    this.update = function () {
        self.scrollerHeight = self.calculateScrollerHeight();
        self.hScrollerWidth = self.calculateScrollerWidth();
        if (self.scrollerHeight / self.scrollContainer.offsetHeight < 1) {
            self.scroller.style.height = self.scrollerHeight + 'px';
        } else {
            self.scroller.style.height = '0px';
        }
        if (self.hScrollerWidth / self.scrollContainer.offsetWidth < 1) {
            self.hScroller.style.width = self.hScrollerWidth + 'px';
        } else {
            self.hScroller.style.width = '0px';
        }
        //sync window position and scroller
        self.moveScroller({deltaY: 0, deltaX: 0});
    };

    this.removeEvents = function () {
        self.scroller.removeEventListener('mousedown', self.startDrag);
        self.hScroller.removeEventListener('mousedown', self.startHDrag);
        this.scrollContentWrapper.removeEventListener('wheel', self.moveScroller);
    };

    this.createScroller();

    // *** Listeners ***
    // Т.к. этот скрипт загружается раньше чем WT, то идем по пути глобального объекта device
    if (window.WT_Device.tablet()) {
        this.scrollContentWrapper.addEventListener('touchmove', self.touchMoveScroller);
        this.scrollContentWrapper.addEventListener('touchstart', self.touchstart);
    } else if (window.WT_Device.desktop()) {
        this.scrollContentWrapper.addEventListener('wheel', self.moveScroller);    
    }
    
};

window.addEventListener('mousemove', function (e) {
    if (!window.scrollbar) return;

    var object = window.scrollbar.object;
    object.scrollBarScroll(e);
});
window.addEventListener('mouseup', function (e) {
    if (!window.scrollbar) return;

    var object = window.scrollbar.object;
    object.stopDrag(e);
});

window.addEventListener('touchmove', function (e) {
    if (!window.scrollbar) return;

    var object = window.scrollbar.object;
    object.scrollBarScroll(e);
});
window.addEventListener('touchend', function (e) {
    if (!window.scrollbar) return;

    var object = window.scrollbar.object;
    object.stopDrag(e);
});