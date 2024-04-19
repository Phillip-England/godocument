



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

        class Zez {
            getState(node, key) {
                return node.getAttribute("zez:" + key).split(" ")
            }
            applyState(node, state) {
                for (let i = 0; i < state.length; i++) {
                    node.classList.add(state[i])
                }
            }
            removeState(node, state) {
                let classListArray = Array.from(node.classList)
                for (let i = 0; i < state.length; i++) {
                    let index = classListArray.indexOf(state[i])
                    if (index > -1) {
                        classListArray.splice(index, 1)
                    }
                
                }
                node.classList = classListArray.join(' ')
            }
            containsState(node, state) {
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
                    this.applyState(nodes[i], this.getState(nodes[i], key))
                }
            }
            toggleState(node, key) {
                let state = this.getState(node, key)
                if (this.containsState(node, state)) {
                    this.removeState(node, state)
                } else {
                    this.applyState(node, this.getState(node, key))
                }
            }
            toggleStateAll(nodes, key) {
                for (let i = 0; i < nodes.length; i++) {
                    let state = this.getState(nodes[i], key)
                    this.toggleState(nodes[i], key, state)
                }
            }
        }

        let zez = new Zez()

        class SiteNav {
            constructor(sitenav, sitenavItems, sitenavDropdowns) {
                this.sitenav = sitenav
                this.sitenavItems = sitenavItems
                this.sitenavDropdowns = sitenavDropdowns
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
                        zez.applyState(item, zez.getState(item, 'active'))
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

        class PageNav {
            constructor(pagenav, pagenavLinks, articleTitles, headerHeight) {
                this.pagenav = pagenav
                this.pagenavLinks = pagenavLinks
                this.articleTitles = articleTitles
                this.headerHeight = headerHeight
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
                        zez.applyState(link, zez.getState(link, 'active'))
                        continue
                    }
                    if (titlePos < this.headerHeight) {
                        zez.removeState(link, zez.getState(link, 'active'))
                        continue
                    }
                    if (!found) {
                        found = true
                        zez.applyState(link, zez.getState(link, 'active'))
                    } else {
                        zez.removeState(link, zez.getState(link, 'active'))
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

            // hooking events and running initializations
            new SiteNav(sitenav, sitenavItems, sitenavDropdowns).hook()
            new PageNav(pagenav, pagenavLinks, articleTitles, header.offsetHeight).hook()


            // init
            Prism.highlightAll();


            // reveal body
            zez.applyState(body, zez.getState(body, 'loaded'))


        }

        eReset(window, 'DOMContentLoaded', onLoad) // initial page load
        eReset(document.getElementsByTagName('body')[0], "htmx:afterOnLoad", onLoad) // after htmx swaps



    }



})();


