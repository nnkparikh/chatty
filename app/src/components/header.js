import React from 'react';

function AppHeader(){
    return (
    <div className="top-bar">
        <div className="random-triangle"> </div>
        <div className="iex-view-text">Chatty</div>
        <div className="links">
            by
            <a className="link-to-me" href="https://neelparikh.net"> nparikh </a>
            <a className="link-to-github" href="https://github.com/nnkparikh/"> github </a>
        </div>
    </div>
    );
}

export default AppHeader;