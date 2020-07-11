'use strict';

// COMP6443{I_FOUND_IT_c9876289e54620cdfa7b4dc5bb66f1f3}

/* eslint browser:true */

const body = document.querySelector('body');
const isDocPage = document
  .querySelector('body.devsite-doc-page') ? true : false;

function highlightActiveNavElement() {
  var elems = document.querySelectorAll('.devsite-section-nav li.devsite-nav-active');
  for (var i = 0; i < elems.length; i++) {
    expandPathAndHighlight(elems[i]);
  }
}

function expandPathAndHighlight(elem) {
  // Walks up the tree from the current element and expands all tree nodes
  var parent = elem.parentElement;
  var parentIsCollapsed = parent.classList.contains('devsite-nav-section-collapsed');
  if (parent.localName === 'ul' && parentIsCollapsed) {
    parent.classList.toggle('devsite-nav-section-collapsed');
    parent.classList.toggle('devsite-nav-section-expanded');
    // Checks if the grandparent is an expandable element
    var grandParent = parent.parentElement;
    var grandParentIsExpandable = grandParent.classList.contains('devsite-nav-item-section-expandable');
    if (grandParent.localName === 'li' && grandParentIsExpandable) {
      var anchor = grandParent.querySelector('a.devsite-nav-toggle');
      anchor.classList.toggle('devsite-nav-toggle-expanded');
      anchor.classList.toggle('devsite-nav-toggle-collapsed');
      expandPathAndHighlight(grandParent);
    }
  }
}

function getCookieValue(name, defaultValue) {
  const value = document.cookie.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)');
  return value ? value.pop() : defaultValue;
}

function initYouTubeVideos() {
  var videoElements = body
    .querySelectorAll('iframe.devsite-embedded-youtube-video');
  videoElements.forEach(function(elem) {
    const videoID = elem.getAttribute('data-video-id');
    if (videoID) {
      let videoURL = 'https://www.youtube.com/embed/' + videoID;
      videoURL += '?autohide=1&amp;showinfo=0&amp;enablejsapi=1';
      elem.src = videoURL;
    }
  });
}

function init() {
  initYouTubeVideos();
  highlightActiveNavElement();
  console.log("%cMMMMWNKOxollcclldxOKWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMW0olllllllllllodxOXWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNx;;dNMMMMMMMMMMM\nMMWKkl:;;lxo:;,;clccokXWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMWk:,:looooooool:;;cONMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNd;;dNMMMMMMMMMMM\nMNkc;;;;oXMWk:ckXNOc;;lONMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMWk:,l0WWWWWWWWNKd:,:OWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNd;;dNMMMMMMMMMMM\nXd:oO0xcxNMMKokWMM0l;;;:xNMMNK000KWMMMMWX0000XWMMMMMMWXK0O00KNWMMMMMMMMMMWXK00O0KXNWMMMMMMMWNK00000KXNWMMMMMWNXK00OO0KXWWMMMMMWO:,lKMMMMMMMMMM0c;;xWMMWNXK00O00KNWMMMWX0KNWNX00O0KXWMMMNd;;dNMMMMMWK00KN\nd;:OWMWOoOXXxlOWW0dok0kl:kWWk:;;;lKMMMMNd;;;;dNMMMNKxoc:;;;;:cokXWMMMMWXkoc:;;;;;:lx0NMMWXOdl:;;;;;:cldONWNOdlc:;;;;;;cokXWMMMWO:;lKMMMMMMMMMNx:;cOWNOdlc:;;;;;:lxKWMNx;:oxl::;;;:cd0WMNd;;dNMMMWKd:;ckX\nc,;oXWMXo:okOkxdlcxNMMXo;lKWk:,,;lKMMMMNd;,;;dNMWKd:;,;;:::;;,;;ckXMMNkc;,;;:cc:;;;:dKMW0l;;;;;:cc:;;;;l0WNd;;;:cclc:;;;;l0WMMWO:;cxOOOOOOOkxl:;lOWMNklldkOO0Oko:;cOWNd;,;:okO00Od:;cOWNd;;dNMWXxc;:dKWM\n:,,;lxkdldKWMMNx:l0WWXx:;c0Wk:,,;lKMMMMNd;;;;dNWKl;;,;lkKXX0x:;;;;dXNk:,;;:xKXNX0xdONWW0c;,,;oOXNXKxodkXWMMXkxOKXNNXOl;;,;dNMMWO:,;;;;;;;;;;;;;:dKWMMWNNMMMMMMMNx;,oXNd;;cOWMMMMMNx;;oXNd;;dNNkc;:o0WMMM\nc,;;;:dOKWMMMMMNOooxdc;,,lKWk:;;;lKMMMMNd;;;;dXNx;,;;l0WMMMMWk:,;;:OKl,;,;xWMMMMMWWMMMNd;;;;oXMMMMMWWWMMMMMWX0kxxddddc;;,;oXMMWk:,ck00000000Oxo:;ckNMWNX0Okkkkkko;;lKNd;;dNMMMMMMWO:;lKNd;;okl;;;dXMMMMM\nd;;;:kWMMMMMMMMMMXx:;,,,:kNWk:;,;lKMMMMNd;,,;dNNd;;;;oXMMMMMM0c;;;:k0l;;;:OWMMMMMMMMMMNd;,;;dNMMMMMMMMMMMMXxc;;,;;;;;;;,,,lXMMWO:,lKMMMMMMMMMMNx:,:OXkl::cclllll:;,lKNd;;dNMMMMMMMO:,lKNd;;;;;:;;ckNMMMM\nXd;,:xNMMMMMMMMMMMXo;;;:xNMWk:,;;c0WMMMKl;;,;dNWk:;;;:kNWMMWXd;;;;l0Xd;;;;oKWMMMWK0XWMWk:;;;cOWMMMWX00KWMNd;;,;lk000Ol;,,,lXMMWO:,lKMMMMMMMMMMWOc;:do;;lOXNNNNNXd;;lKNd;;dNMMMMMMMO:;lKNd;;;cxKOc;:dXWMM\nMNkc;:lxkkOOO0KXK0d:;;,ckNMM0c;;;;lkOOxl;;;;;oNMNx:;,;:lxkkdc;,,;cOWW0l;;,;cdkkxoc;cxKWXd;,,;:oxkkdc:;:d00l;;,;dKXX0d:;,,,lXMMWO:,l0WWWWWWWWWXOl;;lxl;;xNWMMWNKxc;;lKNd;;dNMMMMMMMO:;lKNd;;l0WMWKo;;l0WM\nMMWXkl:;;;;;;;:::;;::;,;:oKWNkc;,;,;;;;;;;;;;dNMMN0o:;;;,;,,,,;cdKWMMWKdc;,,,,;;;;;:l0WMNkl;;;;;;;;,,;;ckXkc;,,;:c:;;:;,;,lXMMWO:,;loooooooolc;;:dKNk:;:ldddol:cc;;lKNd;;dNMMMMMMMO:;lKNd;;dNMMMMXx:;:kN\nMMMMMNKOxdllcclodxOK0d:;:dKWMWKkdllclodkxolllkNMMMMNKkdollcloxOXWMMMMMMWXOxollllodk0NWMMMWN0xdlllllodxOXWMWKkdllcloxOOdlllxXMMM0ollllllllllllodkKWMMWKxdllclodOX0dlxXWOolkNMMMMMMM0olxXWOookNMMMMMNOoloO\nMMMMMMMMMWNNXXNNWMMMMNOoOWMMMMMMWNXXXNWMWNNNNWMMMMMMMMWNNXXNWWMMMMMMMMMMMMWWNXXXNWMMMMMMMMMMMWNXXXNWMMMMMMMMMWNXXXWWMWNNNNNWMMMWNNNNNNNNNNNNNWWMMMMMMMMWNXXXNWMMWNNWWMWNNWMMMMMMMMWNNWWMWNNWMMMMMMMWNNNN", "font-size: 1px")
  console.log("Dreamed of building a safer bank?");
  console.log("We are hiring!");
  console.log("To apply: https://foobar-recruit.quoccabank.com");
}

init();
