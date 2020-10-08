import React from 'react';
import styled from 'styled-components';
import helmLogo from './helm.png';

const NoReleasesContainer = styled.div`
  display: flex;
  height: 100%;
  align-items: center;
`;

function getBackgroundColor(severity) {
    switch (severity) {
        case 'warning':
            return '#8e98d5';
        case 'unknown':
            return '#BA55D3';
        case 'critical':
            return '#9df79a';
        case 'error':
            return '#e7632a';
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

function getHelmLogo() {
    return(
        <img src={helmLogo} width={30} style={{background: 'white', borderRadius: '10px', marginRight: '10px'}} />
    )
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
  cursor: pointer;

  &:hover {
    color: #383d41;
    background-color: #e2e3e5;
    border-color: #d6d8db;
  }
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
  font-weight: ${props => props.isLatest ? 'bold' : 'normal'};
  text-decoration:  ${props => props.isLatest ? 'underline' : 'none'};
 `;

const AppVersionName = styled.span`
  font-size: 0.7em;
  font-weight: ${props => props.isLatest ? 'bold' : 'normal'};
  text-decoration:  ${props => props.isLatest ? 'underline' : 'none'};
 `;


function EventContainer({releases}) {
    return releases.length === 0 ? (
        <NoReleasesContainer>
            <h2>-</h2>
        </NoReleasesContainer>
    ) : (
        <div>
            {
                releases.map(release => {
                    const isLatest = release.type === 'Latest release';
                    const isHelmChart = release.type === 'Helm chart';
                    return (
                        <Alert
                            key={release.repository_name}
                            severity={release.severity}
                            onClick={() => window.open(release.html_url, "_blank")}>
                            <Header severity={release.severity}>
                                <Title>
                                    <strong>{release.repository_name}</strong> - <TagName
                                    isLatest={isLatest}>{release.tag_name}</TagName>
                                    <AppVersionName>{release.app_version} </AppVersionName>
                                </Title>
                                {isHelmChart && getHelmLogo()}
                                <small>{release.days} days</small>
                            </Header>
                        </Alert>
                    )
                })
            }
        </div>

    );
}

export default EventContainer;
