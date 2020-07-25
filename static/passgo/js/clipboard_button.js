(function(){

    const targetButtonSelector = '#copy-clipboard-btn'
    const clipboard = new Clipboard(targetButtonSelector);

    const copyError = function(e) {
        let key;
        if (/Mac/i.test(navigator.userAgent)) {
          key = '&#8984;';
        } else {
          key = 'Ctrl';
        }
        $(e.trigger).attr('title', "Press " + key + "-C to copy" )
                 .tooltip('fixTitle')
                 .tooltip('show');
    };

    const copySuccess = function(e) {
        $(e.trigger).attr('title', 'Copied!')
                 .tooltip('fixTitle')
                 .tooltip('show');
        e.clearSelection();

    };

    clipboard.on('success', copySuccess);
    clipboard.on('error', copyError);

    $(targetButtonSelector).tooltip();

})();