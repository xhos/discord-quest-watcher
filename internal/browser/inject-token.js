const iframe = document.createElement('iframe');
document.body.appendChild(iframe);
iframe.contentWindow.localStorage.token = '"__TOKEN__"';
setTimeout(() => location.reload(), 2000);
