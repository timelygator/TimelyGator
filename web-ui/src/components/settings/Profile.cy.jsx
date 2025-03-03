import React from 'react';
import Profile from './Profile';
import { BrowserRouter as Router } from 'react-router-dom'; // Assuming you're using react-router

describe('<Profile />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(
      <Router>
        <Profile />
      </Router>
    );
  });

  it('clicks the Edit Profile button', () => {
    cy.mount(
      <Router>
        <Profile />
      </Router>
    );

    // Select the Edit Profile button using the data-cy attribute
    cy.get('[data-cy=edit-profile-button]').as('editProfileButton');

    // Verify the button exists and is clickable
    cy.get('@editProfileButton').should('exist').and('be.visible').click();
  });
});