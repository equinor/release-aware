import React from 'react';
import styled from 'styled-components';

const NoEventsContainer = styled.div`
  display: flex;
  height: 100%;
  align-items: center;
`;

function getBackgroundColor(severity) {
  switch (severity) {
    case 'none':
      return '#00b7bf';
    case 'ok':
      return '#00E30F';
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
    case 'none':
      return '#00b7bf';
    case 'ok':
      return '#00E30F';
    case 'warning':
      return '#ffeeba';
    case 'critical':
      return '#f5c6cb';
    default:
      return '#fdfdfe';
  }
}

const EventsContainer = styled.div`
  display: flex;
  height: 100%;
  align-items: center;
`;

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

const Body = styled.div`
  padding: .15rem;
  /*background-color: rgba(255,255,255,.85);
  background-clip: padding-box;
  border-bottom: 1px solid rgba(0,0,0,.05);*/
`;

const Title = styled.span`
  margin-right: auto!important; 
`
function EventContainer({ events }) {
  return events.length === 0 ? (
    <NoEventsContainer>
      <h2>Please specify some repositories to track</h2>
    </NoEventsContainer>
  ) : (
    <div>
          {events.map(event => (
            <EventsContainer>   
            <Alert
              key={event.alertname + event.message}
              severity={event.severity}
              onClick={() => window.open(event.html_url, "_blank")}>
                <Header severity={event.severity}>
                <Title><strong>{event.repository_name}</strong> - <small>{event.tag_name}</small> </Title>
                {false && <small>{new Date(event.published_at).toISOString().slice(0,10)}</small> }
                <small>{event.days} days</small>
              </Header>
              <Body></Body>
            </Alert>                         
            </EventsContainer>
          ))}       
    </div>
   
  );
}

export default EventContainer;
