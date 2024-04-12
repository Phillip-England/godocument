



(() => {

    if (document.getElementsByTagName('html')[0].getAttribute('loaded') == 'false') {

        document.getElementsByTagName('html')[0].setAttribute('loaded', 'true')
        sessionStorage.setItem('sitenav', '')

        function climbTreeUntil(element, callback) {
            let parent = element.parentElement;
            if (callback(parent)) {
                return parent
            }
            return climbTreeUntil(parent, callback)
        }
    
        class SiteNav {
            constructor() {
                this.nav = document.querySelector('#sitenav');
                this.pullNavFromSession()
                this.items = document.querySelectorAll('.sitenav-item');
                this.buttons = document.querySelectorAll('.sitenav-dropdown-button');
                this.overlay = document.querySelector('#mobile-nav-overlay');
                this.overlay.removeEventListener('click', this.toggleMobileNav.bind(this));
                this.overlay.addEventListener('click', this.toggleMobileNav.bind(this));
                this.buttons.forEach((button) => {
                    button.classList.remove('sitenav-dropdown-button-active')
                    button.removeEventListener('click', this.eToggleDropdown.bind(this))
                    button.addEventListener('click', this.eToggleDropdown.bind(this))
                })
                this.setActiveLink()
            }
            pullNavFromSession() {
                let nav = sessionStorage.getItem('sitenav')
                if (nav) {
                    this.nav.innerHTML = nav
                }
            }
            setActiveLink() {
                this.items.forEach((item) => {
                    if (item.getAttribute('href') == window.location.pathname) {
                        item.classList.add('sitenav-item-active');
                        this.openActiveLinkParents(item)
                    } else {
                        item.classList.remove('sitenav-item-active');
                    }
                });
            }
            openActiveLinkParents(element) {
                let parent = element.parentElement;
                if (parent.id == 'sitenav') {
                    return
                }
                if (parent.classList.contains('sitenav-dropdown')) {
                    let dropdownButton = parent.firstChild
                    dropdownButton.classList.add('sitenav-dropdown-button-active');
                    let caret = dropdownButton.querySelector('.dropdown-caret')
                    caret.classList.add('dropdown-caret-active')
                    let hiddenSection = parent.querySelector('.sitenav-dropdown-children');
                    hiddenSection.classList.remove('hidden')
                }
                this.openActiveLinkParents(parent)
            }
            eToggleDropdown(e) {
                let dropdown = climbTreeUntil(e.target, (element) => {
                    return element.classList.contains('sitenav-dropdown')
                });
                let childDropdowns = dropdown.querySelectorAll('.sitenav-dropdown');
                let dropdownButton = dropdown.firstChild
                let dropdownCaret = dropdownButton.querySelector('.dropdown-caret')
                if (dropdownCaret.classList.contains('dropdown-caret-active')) {
                    this.closeDropdown(dropdown)
                    for (let i = 0; i < childDropdowns.length; i++) {
                        this.closeDropdown(childDropdowns[i])
                    }
                } else {
                    this.openDropdown(dropdown)
                }
                sessionStorage.setItem('sitenav', this.nav.innerHTML)
            }
            openDropdown(element) {
                let dropdownButton = element.firstChild
                let caret = dropdownButton.querySelector('.dropdown-caret')
                let hiddenSection = element.querySelector('.sitenav-dropdown-children');
                caret.classList.add('dropdown-caret-active')
                hiddenSection.classList.remove('hidden')
            }
            closeDropdown(element) {
                let dropdownButton = element.firstChild
                let caret = dropdownButton.querySelector('.dropdown-caret')
                let hiddenSection = element.querySelector('.sitenav-dropdown-children');
                caret.classList.remove('dropdown-caret-active')
                hiddenSection.classList.add('hidden')
            }
            toggleMobileNav() {
                if (this.nav.classList.contains('sitenav-active')) {
                    this.overlay.classList.remove('mobile-nav-overlay-active')
                    this.nav.classList.remove('sitenav-active');
                    return
                }
                this.overlay.classList.add('mobile-nav-overlay-active')
                this.nav.classList.add('sitenav-active');
            }
        }

        class Header {
            constructor(sitenav) {
                this.sitenav = sitenav;
                this.header = document.querySelector('#header');
                this.bars = header.querySelector('.header-bars');
                this.bars.removeEventListener('click', this.toggleMobileNav.bind(this));
                this.bars.addEventListener('click', this.toggleMobileNav.bind(this));
                this.headerHeight = this.header.offsetHeight;
            }
            toggleMobileNav() {
                this.sitenav.toggleMobileNav();
            }        
        }

        class Theme {
            constructor() {
                this.sunIcons = document.querySelectorAll('.sun-icon');
                for (let i = 0; i < this.sunIcons.length; i++) {
                    this.sunIcons[i].removeEventListener('click', this.toggleTheme.bind(this));
                    this.sunIcons[i].addEventListener('click', this.toggleTheme.bind(this));
                }
                this.moonIcons = document.querySelectorAll('.moon-icon');
                for (let i = 0; i < this.moonIcons.length; i++) {
                    this.moonIcons[i].removeEventListener('click', this.toggleTheme.bind(this));
                    this.moonIcons[i].addEventListener('click', this.toggleTheme.bind(this));
                }
                this.theme = localStorage.getItem('theme');
                if (!this.theme) {
                    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
                        this.theme = 'dark';
                    } else {
                        this.theme = 'light';
                    }
                    localStorage.setItem('theme', this.theme);
                }
                this.applyTheme();
            }
            applyTheme() {
                document.documentElement.setAttribute('class', this.theme);
            }
            toggleTheme() {
                if (this.theme == 'light') {
                    this.theme = 'dark';
                } else {
                    this.theme = 'light';
                }
                localStorage.setItem('theme', this.theme);
                this.applyTheme();
            }
        }

        class Article {
            constructor(headerHeight) {
                this.article = document.getElementById('article')
                this.h2 = this.article.querySelectorAll('h2')
                this.h3 = this.article.querySelectorAll('h3')
                this.h4 = this.article.querySelectorAll('h4')
                this.h5 = this.article.querySelectorAll('h5')
                this.h6 = this.article.querySelectorAll('h6')
                this.headerHeight = headerHeight
                this.headers = []
                this.gatherHeaders()
                this.setActiveHeader()
                this.article.removeEventListener('scroll', this.setActiveHeader.bind(this))
                this.article.addEventListener('scroll', this.setActiveHeader.bind(this))
                
            }
            gatherHeaders() {
                for (let i = 0; i < this.h2.length; i++) {
                    this.headers.push(this.h2[i])
                }
                for (let i = 0; i < this.h3.length; i++) {
                    this.headers.push(this.h3[i])
                }
                for (let i = 0; i < this.h4.length; i++) {
                    this.headers.push(this.h4[i])
                }
                for (let i = 0; i < this.h5.length; i++) {
                    this.headers.push(this.h5[i])
                }
                for (let i = 0; i < this.h6.length; i++) {
                    this.headers.push(this.h6[i])
                }
            }
            setActiveHeader() {
                let found = false
                for (let i = 0; i < this.headers.length; i++) {
                    let header = this.headers[i]
                    let headerTop = this.headers[i].getBoundingClientRect().top
                    if (headerTop < this.headerHeight) {
                        header.classList.remove('article-header-active')
                        continue
                    }
                    if (found) {
                        header.classList.remove('article-header-active')
                        continue
                    }
                    header.classList.add('article-header-active')
                    found = true
                }
            }
        }

        class PageNav {
            constructor() {
                this.nav = document.querySelector('#pagenav');
            }
        }

        function onLoad() {
            let sitenav = new SiteNav();
            let header = new Header(sitenav);
            let theme = new Theme()
            let article = new Article(header.headerHeight)
            console.log(article)
            Prism.highlightAll();
        }

        


    
    }
    
    window.removeEventListener('DOMContentLoaded', onLoad);
    window.addEventListener('DOMContentLoaded', onLoad);

    document.getElementsByTagName('body')[0].removeEventListener("htmx:afterOnLoad", onLoad)
    document.getElementsByTagName('body')[0].addEventListener("htmx:afterOnLoad", onLoad)



})();


