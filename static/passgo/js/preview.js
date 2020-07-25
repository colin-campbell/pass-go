(function () {

  $('#revealSecret').click(function () {
    const form = $('<form/>')
      .attr('id', 'revealSecretForm')
      .attr('method', 'post');
    form.appendTo($('body'));
    form.submit();
  });
})();