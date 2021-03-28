import '../style/404.scss';


window.addEventListener('load', () => {
    let container = document.getElementsByClassName("container")[0];
    container.addEventListener('click', () => {
        location.href = '/';
    });
});