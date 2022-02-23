function fscroll(){
    if (window.location.href.toString().indexOf("scroll") !== -1) {
        let sblock = document.getElementsByClassName("block")[1]
        let scr = sblock.offsetTop
        window.scroll({top: scr, left: 0, behavior: "auto"})
    }
}

fscroll()