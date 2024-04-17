



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
                    if (item.getAttribute('href') == window.location.pathname || item.getAttribute('href') == window.location.href) {
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
                    let dropdownButton = parent.children[0]
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
                let dropdownButton = dropdown.children[0]
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
                let dropdownButton = element.children[0]
                let caret = dropdownButton.querySelector('.dropdown-caret')
                let hiddenSection = element.querySelector('.sitenav-dropdown-children');
                caret.classList.add('dropdown-caret-active')
                hiddenSection.classList.remove('hidden')
            }
            closeDropdown(element) {
                let dropdownButton = element.children[0]
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
            constructor(headerHeight, pagenavLinks) {
                this.article = document.getElementById('article')
                this.contentWrapper = document.getElementById('content-wrapper')
                this.pagenavLinks = pagenavLinks
                this.h2 = this.article.querySelectorAll('h2')
                this.h3 = this.article.querySelectorAll('h3')
                this.h4 = this.article.querySelectorAll('h4')
                this.h5 = this.article.querySelectorAll('h5')
                this.h6 = this.article.querySelectorAll('h6')
                this.headerHeight = headerHeight
                this.headers = []
                this.scrollTimeout = null
                this.gatherHeaders()
                this.setActivePagenavLink()
                this.setBlankArticleLinks()
                this.article.removeEventListener('scroll', this.articleScrollEvent.bind(this))
                this.article.addEventListener('scroll', this.articleScrollEvent.bind(this))
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
            setActivePagenavLink() {
                let found = false
                for (let i = 0; i < this.headers.length; i++) {
                    let header = this.headers[i]
                    let headerTop = this.headers[i].getBoundingClientRect().top
                    if (i == this.headers.length - 1 && headerTop < this.headerHeight) {
                        this.pagenavLinks[i].classList.add('pagenav-link-active')
                        continue
                    }
                    if (headerTop < this.headerHeight) {
                        this.pagenavLinks[i].classList.remove('pagenav-link-active')
                        continue
                    }
                    if (found) {
                        this.pagenavLinks[i].classList.remove('pagenav-link-active')
                        continue
                    }
                    this.pagenavLinks[i].classList.add('pagenav-link-active')
                    found = true

                }
            }
            articleScrollEvent() {
                if (this.scrollTimeout !== null) {
                    clearTimeout(this.scrollTimeout);
                }
                this.scrollTimeout = setTimeout(() => {
                    this.setActivePagenavLink();
                }, 100);
            }
            setBlankArticleLinks() {
                let links = this.contentWrapper.querySelectorAll('a');
                for (let i = 0; i < links.length; i++) {
                    let link = links[i];
                    let linkTarget = link.getAttribute('target');
                    if (!linkTarget) {
                        link.setAttribute('target', '_blank');
                    }
                }
            }
            scrollToArticleHeader(header) {
                let headerTop = header.getBoundingClientRect().top;
                let articleTop = this.article.getBoundingClientRect().top;
                console.log(headerTop, articleTop)
            }
        }

        class PageNav {
            constructor() {
                this.nav = document.querySelector('#pagenav');
                this.links = this.nav.querySelectorAll('a');
                // remove any punctuation from link hrefs
                for (let i = 0; i < this.links.length; i++) {
                    let link = this.links[i];
                    let href = link.getAttribute('href');
                    href = href.replace("?", "")
                    href = href.replace("!", "")
                    href = href.replace(".", "")
                    href = href.replace("/", "")
                    href = href.replace(".", "")
                    link.setAttribute('href', href);
                }
            }
        }

        class CustomContentElements {
            constructor() {
                this.contentImportant = document.querySelectorAll('.content-important');
                this.contentWarning = document.querySelectorAll('.content-warning');
                this.insertExclamationPoints()
                this.insertExplodingFaces()
            }
            insertExclamationPoints() {
                for (let i = 0; i < this.contentImportant.length; i++) {
                    let parentP = this.contentImportant[i].parentElement;
                    let exclamation = document.createElement('div')
                    exclamation.classList.add('content-important-exclamation')
                    exclamation.innerHTML = `
                        <svg aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="none" viewBox="0 0 24 24">
                            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 13V8m0 8h.01M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
                        </svg>
                        <p>Important</p>
                    `
                    parentP.insertBefore(exclamation, this.contentImportant[i])
                    
                }
            }
            insertExplodingFaces() {
                for (let i = 0; i < this.contentWarning.length; i++) {
                    let parentP = this.contentWarning[i].parentElement;
                    let explodingFace = document.createElement('div')
                    explodingFace.classList.add('content-warning-explode')
                    explodingFace.innerHTML = `
                        <svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                            <path fill="currentColor" d="M12 17a2 2 0 0 1 2 2h-4a2 2 0 0 1 2-2Z"/>
                            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.815 9H16.5a2 2 0 1 0-1.03-3.707A1.87 1.87 0 0 0 15.5 5 1.992 1.992 0 0 0 12 3.69 1.992 1.992 0 0 0 8.5 5c.002.098.012.196.03.293A2 2 0 1 0 7.5 9h3.388m2.927-.985v3.604M10.228 9v2.574M15 16h.01M9 16h.01m11.962-4.426a1.805 1.805 0 0 1-1.74 1.326 1.893 1.893 0 0 1-1.811-1.326 1.9 1.9 0 0 1-3.621 0 1.8 1.8 0 0 1-1.749 1.326 1.98 1.98 0 0 1-1.87-1.326A1.763 1.763 0 0 1 8.46 12.9a2.035 2.035 0 0 1-1.905-1.326A1.9 1.9 0 0 1 4.74 12.9 1.805 1.805 0 0 1 3 11.574V12a9 9 0 0 0 18 0l-.028-.426Z"/>
                        </svg>
                        <p>Warning</p>
                    `
                    parentP.insertBefore(explodingFace, this.contentWarning[i])
                }
            }
        }

        function onLoad() {
            let sitenav = new SiteNav();
            let header = new Header(sitenav);
            let theme = new Theme()
            let pagenav = new PageNav()
            let article = new Article(header.headerHeight, pagenav.links)
            let customContentElements = new CustomContentElements()
            Prism.highlightAll();
        }

        


    
    }
    
    window.removeEventListener('DOMContentLoaded', onLoad);
    window.addEventListener('DOMContentLoaded', onLoad);

    document.getElementsByTagName('body')[0].removeEventListener("htmx:afterOnLoad", onLoad)
    document.getElementsByTagName('body')[0].addEventListener("htmx:afterOnLoad", onLoad)



})();


