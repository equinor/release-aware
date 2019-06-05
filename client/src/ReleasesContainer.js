import React from 'react';
import styled from 'styled-components';

const NoReleasesContainer = styled.div`
  display: flex;
  height: 100%;
  align-items: center;
`;

function getBackgroundColor(severity) {
    switch (severity) {
        case 'warning':
            return '#fff3cd';
        case 'critical':
            return '#f8d7da';
        default:
            return '#fefefe';
    }
}

function getBorderColor(severity) {
    switch (severity) {
        case 'warning':
            return '#ffeeba';
        case 'critical':
            return '#f5c6cb';
        default:
            return '#fdfdfe';
    }
}


const Alert = styled.div`
  font-size: 1em;
  width: 450px;
  overflow: hidden;
  background-clip: padding-box;
  border: 1px solid rgba(0,0,0,.1);
  box-shadow: 0 0.25rem 0.75rem rgba(0,0,0,.1);
  -webkit-backdrop-filter: blur(10px);
  backdrop-filter: blur(10px);
  border-radius: .25rem;
  border-color: ${props => getBorderColor(props.severity)};
  background: ${props => getBackgroundColor(props.severity)};  
  margin-bottom: 10px;
  margin-left: 10px;
`;

const Header = styled.div`
  display: flex;
  -ms-flex-align: center;
  align-items: center;
  padding: .25rem .55rem;
  color: black;
`;

const Title = styled.span`
  margin-right: auto!important; 
`;

const TagName = styled.span`
  font-size: 0.7em;
  font-weight: ${props => props.isLatest ? 'bold' : 'normal' };  
  text-decoration:  ${props => props.isLatest ? 'underline' : 'none' };
 `;


function EventContainer({events}) {
    return events.length === 0 ? (
        <NoReleasesContainer>
            <h2>Please specify some repositories to track</h2>
        </NoReleasesContainer>
    ) : (
        <div>
            {events.map(event => (
                <Alert
                    key={event.alertname + event.message}
                    severity={event.severity}
                    onClick={() => window.open(event.html_url, "_blank")}>
                    <Header severity={event.severity}>
                        <Title>
                            <strong>{event.repository_name}</strong> - <TagName isLatest={event.type === 'Latest release'}>{event.tag_name}</TagName>
                        </Title>
                        <small>{event.days} days</small>
                    </Header>
                </Alert>
            ))}
        </div>

    );
}

export default EventContainer;
