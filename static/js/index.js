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
                let exit = callback(node)
                if (exit == true) {
                    return
                }
                climbTreeUntil(node.parentNode, stopNode, callback)
            }
        }

        function getLoadState() {
            return document.getElementsByTagName('html')[0].getAttribute('loaded')
        }

        function doOnce(callback) {
            if (getLoadState() == 'false') {
                callback()
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
                        return true
                    }
                })
                let hiddenChildren = qs(dropdown, 'ul')
                let caret = qs(dropdown, 'div')
                zez.toggleStateAll([caret, hiddenChildren], 'active')
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
            constructor(pagenav, pagenavLinks, articleTitles) {
                this.pagenav = pagenav
                this.pagenavLinks = pagenavLinks
                this.articleTitles = articleTitles
                this.windowTimeout = null
                this.bufferZone = 200
                this.activeLink = null
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
                for (let i = 0; i < this.articleTitles.length; i++) {
                    let link = this.pagenavLinks[i]
                    let nextLink = this.pagenavLinks[i + 1]
                    let title = this.articleTitles[i]
                    let titlePos = title.getBoundingClientRect().top
                    let nextTitle = this.articleTitles[i + 1]
                    let nextTitlePos = nextTitle ? nextTitle.getBoundingClientRect().top : 0
                    if (i == 0 && titlePos > 0) {
                        this.activeLink = link
                        break
                    }
                    if (i == this.articleTitles.length - 1 && titlePos < 0) {
                        this.activeLink = link
                        break
                    }
                    if (titlePos < 0 && nextTitlePos > 0) {
                        if (nextTitlePos < this.bufferZone) {
                            this.activeLink = nextLink
                            continue
                        }
                        this.activeLink = link
                    }
                }
                for (let i = 0; i < this.pagenavLinks.length; i++) {
                    let link = this.pagenavLinks[i]
                    if (link == this.activeLink) {
                        zez.applyState(link, 'active')
                    } else {
                        zez.removeState(link, 'active')
                    }
                }
            }
            handleWindowScroll() {
                clearTimeout(this.windowTimeout);
                this.windowTimeout = setTimeout(() => {
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

        class CustomComponents {
            initComponent(node) {
                node.parentElement.replaceWith(node)
                let text = node.innerHTML
                text = this.replaceBackticksWithCodeTags(text)
                return text
            }
            replaceBackticksWithCodeTags(text) {
                for (let i = 0; i < text.length; i++) {
                    if (text[i] == '`') {
                        text = text.slice(0, i) + '<code>' + text.slice(i + 1)
                        i++
                        while (i < text.length && text[i] != '`') {
                            i++
                        }
                        text = text.slice(0, i) + '</code>' + text.slice(i + 1)
                    }
                }
                return text
            }
            registerComponent(className, htmlContent) {
                let elements = document.getElementsByClassName(className)
                for (let i = 0; i < elements.length; i++) {
                    let node = elements[i]
                    let text = this.initComponent(node)
                    node.innerHTML = htmlContent.replace('{text}', text)
                }
            }
        }

// ==============================================================================

class MdImportant {
    constructor(customComponents) {
        customComponents.registerComponent("md-important", `
            <div class='bg-[var(--md-bg-color)] dark:bg-[var(--dark-md-bg-color)] p-4 rounded-md border-l-4 border-[var(--md-important-border-color)] dark:border-[var(--dark-md-important-border-color)] flex flex-col gap-2'>
                <span class='flex flex-row item-center gap-2 dark:text-[var(--dark-md-important-text-color)] text-[var(--md-important-text-color)]'>
                    <span class='flex items-center dark:text-[var(--dark-md-important-text-color)] text-[var(--md-important-text-color)]'>
                        <svg class="h-6 w-6" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 13V8m0 8h.01M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
                        </svg>
                    </span>
                    <p class='font-bold'>Important</p>                   
                </span>
                <p class='custom-inline-code'>{text}</p>
            </div>
        `)
    }
}

// ==============================================================================

class MdWarning {
    constructor(customComponents) {
        customComponents.registerComponent("md-warning", `
            <div class='bg-[var(--md-bg-color)] dark:bg-[var(--dark-md-bg-color)] p-4 rounded-md border-l-4 border-[var(--md-warning-border-color)] dark:border-[var(--dark-md-warning-border-color)] flex flex-col gap-2'>
                <span class='flex flex-row item-center gap-2 dark:text-[var(--dark-md-warning-text-color)] text-[var(--md-warning-text-color)]'>
                    <span class='flex items-center dark:text-[var(--dark-md-warning-text-color)] text-[var(--md-warning-text-color)]'>
                    <svg class="w-6 h-6" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                        <path fill="currentColor" d="M12 17a2 2 0 0 1 2 2h-4a2 2 0 0 1 2-2Z"/>
                        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.815 9H16.5a2 2 0 1 0-1.03-3.707A1.87 1.87 0 0 0 15.5 5 1.992 1.992 0 0 0 12 3.69 1.992 1.992 0 0 0 8.5 5c.002.098.012.196.03.293A2 2 0 1 0 7.5 9h3.388m2.927-.985v3.604M10.228 9v2.574M15 16h.01M9 16h.01m11.962-4.426a1.805 1.805 0 0 1-1.74 1.326 1.893 1.893 0 0 1-1.811-1.326 1.9 1.9 0 0 1-3.621 0 1.8 1.8 0 0 1-1.749 1.326 1.98 1.98 0 0 1-1.87-1.326A1.763 1.763 0 0 1 8.46 12.9a2.035 2.035 0 0 1-1.905-1.326A1.9 1.9 0 0 1 4.74 12.9 1.805 1.805 0 0 1 3 11.574V12a9 9 0 0 0 18 0l-.028-.426Z"/>
                    </svg>
                
                    </span>
                    <p class='font-bold'>Warning</p>                   
                </span>
                <p class='custom-inline-code'>{text}</p>
            </div>            
        `)
    }
}

// ==============================================================================

class MdCorrect {
    constructor(customComponents) {
        customComponents.registerComponent("md-correct", `
            <div class='bg-[var(--md-bg-color)] dark:bg-[var(--dark-md-bg-color)] p-4 rounded-md border-l-4 border-[var(--md-correct-border-color)] dark:border-[var(--dark-md-correct-border-color)] flex flex-col gap-2'>
                <span class='flex flex-row item-center gap-2 dark:text-[var(--dark-md-correct-text-color)] text-[var(--md-correct-text-color)]'>
                    <span class='flex items-center dark:text-[var(--dark-md-correct-text-color)] text-[var(--md-correct-text-color)]'>
                    <svg class="w-6 h-6" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.5 11.5 11 14l4-4m6 2a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
                    </svg>
                </span>
                    <p class='font-bold'>Correct</p>                   
                </span>
                <p class='custom-inline-code'>{text}</p>
            </div>
        `)
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
            new PageNav(pagenav, pagenavLinks, articleTitles)
            new Header(headerBars, overlay, sitenav)
            new Theme(sunIcons, moonIcons, htmlDocument)

            // defining custom component
            let customComponents = new CustomComponents()
            new MdImportant(customComponents)
            new MdWarning(customComponents)
            new MdCorrect(customComponents)

            // init
            Prism.highlightAll();

            // reveal body
            zez.applyState(body, 'loaded')

            // set loaded attribute
            document.getElementsByTagName('html')[0].setAttribute('loaded', 'true')

        }

        eReset(window, 'htmx:beforeHistorySave', () => {
            document.getElementsByTagName('html')[0].setAttribute('loaded', 'false')
        }) // initial page load
        eReset(window, 'DOMContentLoaded', onLoad) // initial page load
        eReset(document.getElementsByTagName('body')[0], "htmx:afterOnLoad", onLoad) // after htmx swaps

    }

})();


