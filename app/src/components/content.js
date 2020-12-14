import React from 'react';
import IEXToolbar from './toolbar';

function AppContent() {
    return (
    <div className="content-wrapper">
        <div className="main-content">
            <IEXToolbar />
            <div className="company-details"></div>
        </div>
    </div>
    );
}

export default AppContent;