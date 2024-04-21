/*
README
Here is how this script works.
We are using hx-boost (https://htmx.org/attributes/hx-boost/) to load pages without refreshing the entire page.
This means the <head> of our document is never reloaded.
The <head> is loaded on the initial page load, and then after navigating, the new HTML content is loaded into the <body> of the document.
At the bottom of this file, we are referencing the htmx event "htmx:afterOnLoad" to run the script after the new content is loaded. 
So, on the initial page load, we use the DOMContentLoaded event to run the script, and then after navigating, we use htmx:afterOnLoad to run the script.
This changes the way we have to hook events
Everytime you navigate, you have to detach and re-hook events to elements
Failing to detach events will cause the event to be fired multiple times
The utility function eReset(node, eventType, callback) is used to detach and re-hook events
*/


(() => {

    if (document.getElementsByTagName('html')[0].getAttribute('loaded') == 'false') {

        document.getElementsByTagName('html')[0].setAttribute('loaded', 'true')
        sessionStorage.setItem('sitenav', '')

        function qs(root, selector) {
            if (!root) {
                console.error('Root is not defined in qs()')
            }
            return root.querySelector(selector)
        }

        function qsa(root, selector) {
            if (!root) {
                console.error('Root is not defined in qsa()')
            }
            return root.querySelectorAll(selector)
        }

        function climbTreeUntil(node, stopNode, callback) {
            if (node) {                
                if (node == stopNode) {
                    return
                }
                callback(node)
                climbTreeUntil(node.parentNode, stopNode, callback)
            }
        }

        function eReset(node, eventType, callback) {
            node.removeEventListener(eventType, callback)
            node.addEventListener(eventType, callback)
        }

// ==============================================================================

        class Zez {
            getState(node, key) {
                return node.getAttribute("zez:" + key).split(" ")
            }
            applyState(node, stateKey) {
                let state = this.getState(node, stateKey)
                for (let i = 0; i < state.length; i++) {
                    node.classList.add(state[i])
                }
            }
            removeState(node, key) {
                let state = this.getState(node, key)
                let classListArray = Array.from(node.classList)
                for (let i = 0; i < state.length; i++) {
                    let index = classListArray.indexOf(state[i])
                    if (index > -1) {
                        classListArray.splice(index, 1)
                    }
                
                }
                node.classList = classListArray.join(' ')
            }
            containsState(node, key) {
                let state = this.getState(node, key)
                let classListArray = Array.from(node.classList)
                for (let i = 0; i < state.length; i++) {
                    if (classListArray.includes(state[i])) {
                        return true
                    }
                }
                return false
            }
            applyStateAll(nodes, key) {
                for (let i = 0; i < nodes.length; i++) {
                    this.applyState(nodes[i], key)
                }
            }
            toggleState(node, key) {
                let containsState = this.containsState(node, key)
                if (containsState) {
                    this.removeState(node, key)
                } else {
                    this.applyState(node, key)
                }
            }
            swapStates(node, key1, key2) {
                if (this.containsState(node, key1)) {
                    this.enforceState(node, key2, key1)
                } else {
                    this.enforceState(node, key1, key2)
                }
            }
            toggleStateAll(nodes, key) {
                for (let i = 0; i < nodes.length; i++) {
                    this.toggleState(nodes[i], key)
                }
            }
            enforceState(node, keyToApply, keyToRemove) {
                this.applyState(node, keyToApply)
                this.removeState(node, keyToRemove)
            }
        }

        let zez = new Zez()

// ==============================================================================

        class SiteNav {
            constructor(sitenav, sitenavItems, sitenavDropdowns) {
                this.sitenav = sitenav
                this.sitenavItems = sitenavItems
                this.sitenavDropdowns = sitenavDropdowns
                this.hook()
            }
            hook() {
                this.setActiveNavItem()
                for (let i = 0; i < this.sitenavDropdowns.length; i++) {
                    eReset(qs(this.sitenavDropdowns[i], 'button'), "click", this.toggleDropdown.bind(this))
                }
            }
            toggleDropdown(e) {
                let dropdown = null
                climbTreeUntil(e.target, this.sitenav, (node) => {
                    if (node.tagName == 'LI') {
                        dropdown = node
                    }
                })
                let hiddenChildren = qs(dropdown, 'ul')
                let caret = qs(dropdown, 'div')
                let summary = qs(dropdown, 'summary')
                zez.toggleStateAll([summary, caret, hiddenChildren], 'active')
            }
            setActiveNavItem() {
                for (let i = 0; i < this.sitenavItems.length; i++) {
                    let item = this.sitenavItems[i]
                    let href = item.getAttribute('href')
                    if (href == window.location.pathname || href == window.location.href) {
                        zez.applyState(item, 'active')
                        climbTreeUntil(item, this.sitenav, (node) => {
                            if (node.classList.contains('dropdown')) {
                                let hiddenChildren = qs(node, 'ul')
                                let caret = qs(node, 'div')
                                let summary = qs(node, 'summary')
                                zez.applyStateAll([summary, caret, hiddenChildren], 'active')
                            }
                        })
                    }
                }
            }
        }

// ==============================================================================

        class PageNav {
            constructor(pagenav, pagenavLinks, articleTitles, headerHeight) {
                this.pagenav = pagenav
                this.pagenavLinks = pagenavLinks
                this.articleTitles = articleTitles
                this.headerHeight = headerHeight
                this.hook()
            }
            hook() {
                this.setActivePageNavItem()
                eReset(window, "scroll", this.handleWindowScroll.bind(this))
            }
            setActivePageNavItem() {
                if (this.pagenavLinks.length == 0 || this.articleTitles.length == 0) {
                    return
                }
                let found = false
                for (let i = 0; i < this.articleTitles.length; i++) {
                    let link = this.pagenavLinks[i]
                    let title = this.articleTitles[i]
                    let titlePos = title.getBoundingClientRect().top
                    if (!found && i == this.articleTitles.length - 1) {
                        zez.applyState(link, 'active')
                        continue
                    }
                    if (titlePos < this.headerHeight) {
                        zez.removeState(link, 'active')
                        continue
                    }
                    if (!found) {
                        found = true
                        zez.applyState(link, 'active')
                    } else {
                        zez.removeState(link, 'active')
                    }
                }
            }
            handleWindowScroll() {
                let windowTimeout;
                clearTimeout(windowTimeout);
                windowTimeout = setTimeout(() => {
                    this.setActivePageNavItem()
                }, 100);
            }
        }

// ==============================================================================

class Header {
    constructor(headerBars, overlay, sitenav) {
        this.headerBars = headerBars
        this.overlay = overlay
        this.sitenav = sitenav
        this.hook()
    }
    hook() {
        eReset(this.headerBars, "click", this.toggleMobileNav.bind(this))
        eReset(this.overlay, "click", this.toggleMobileNav.bind(this))
    }
    toggleMobileNav() {
        zez.toggleState(this.overlay, 'active')
        zez.toggleState(this.sitenav, 'active')
    }
}

// ==============================================================================

class Theme {
    constructor(sunIcons, moonIcons, htmlDocument) {
        this.sunIcons = sunIcons
        this.moonIcons = moonIcons
        this.htmlDocument = htmlDocument
        this.hook()
    }
    hook() {
        this.initTheme()
        for (let i = 0; i < this.sunIcons.length; i++) {
            eReset(this.sunIcons[i], "click", this.toggleTheme.bind(this))
        }
        for (let i = 0; i < this.moonIcons.length; i++) {
            eReset(this.moonIcons[i], "click", this.toggleTheme.bind(this))
        }
    }
    initTheme() {
        let theme = localStorage.getItem('theme')
        if (theme) {
            if (theme == 'dark') {
                zez.enforceState(this.htmlDocument, 'dark', 'light')
                return
            }
            zez.enforceState(this.htmlDocument, 'light', 'dark')
            return
        }
        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            zez.enforceState(this.htmlDocument, 'dark', 'light')
        } else {
            zez.enforceState(this.htmlDocument, 'light', 'dark')
        }
    }
    toggleTheme() {
        zez.swapStates(this.htmlDocument, 'dark', 'light')
        if (zez.containsState(this.htmlDocument, 'dark')) {
            localStorage.setItem('theme', 'dark')
            return
        }
        localStorage.setItem('theme', 'light')
    }
}

// ==============================================================================



        function onLoad() {

            // elements
            const body = qs(document, 'body')
            const sitenav = qs(document, '#sitenav')
            const sitenavItems = qsa(sitenav, '.item')
            const sitenavDropdowns = qsa(sitenav, '.dropdown')
            const pagenav = qs(document, '#pagenav')
            const pagenavLinks = qsa(pagenav, 'a')
            const article = qs(document, '#article')
            const articleTitles = qsa(article, 'h2, h3, h4, h5, h6')
            const header = qs(document, '#header')
            const headerBars = qs(header, '#bars')
            const overlay = qs(document, '#overlay')
            const sunIcons = qsa(document, '.sun')
            const moonIcons = qsa(document, '.moon')
            const htmlDocument = qs(document, 'html')

            // hooking events and running initializations
            window.scrollTo(0, 0, { behavior: 'auto' })
            new SiteNav(sitenav, sitenavItems, sitenavDropdowns, header, overlay)
            new PageNav(pagenav, pagenavLinks, articleTitles, header.offsetHeight)
            new Header(headerBars, overlay, sitenav)
            new Theme(sunIcons, moonIcons, htmlDocument)

            // init
            Prism.highlightAll();

            // reveal body
            zez.applyState(body, 'loaded')


        }

        eReset(window, 'DOMContentLoaded', onLoad) // initial page load
        eReset(document.getElementsByTagName('body')[0], "htmx:afterOnLoad", onLoad) // after htmx swaps

    }

})();


