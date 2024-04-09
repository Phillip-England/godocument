



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
    
        function getRawHTML(element) {
            var wrapper = document.createElement('div');
            wrapper.appendChild(element.cloneNode(true));
            return wrapper.innerHTML;
        }
    
        class SiteNav {
            constructor() {
                this.nav = document.querySelector('#sitenav');
                this.pullNavFromSession()
                this.items = document.querySelectorAll('.sitenav-item');
                this.buttons = document.querySelectorAll('.sitenav-dropdown-button');
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
                // dropdownButton.classList.add('sitenav-dropdown-active');
                caret.classList.add('dropdown-caret-active')
                hiddenSection.classList.remove('hidden')
            }
            closeDropdown(element) {
                let dropdownButton = element.firstChild
                let caret = dropdownButton.querySelector('.dropdown-caret')
                let hiddenSection = element.querySelector('.sitenav-dropdown-children');
                // dropdownButton.classList.remove('sitenav-dropdown-active');
                caret.classList.remove('dropdown-caret-active')
                hiddenSection.classList.add('hidden')
            }

        }
    
        function onLoad() {
            new SiteNav();
        }

    
    }
    
    window.removeEventListener('DOMContentLoaded', onLoad);
    window.addEventListener('DOMContentLoaded', onLoad);

    document.getElementsByTagName('body')[0].removeEventListener("htmx:afterOnLoad", onLoad)
    document.getElementsByTagName('body')[0].addEventListener("htmx:afterOnLoad", onLoad)



})();


