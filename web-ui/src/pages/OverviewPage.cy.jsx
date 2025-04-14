import React from 'react';
import OverviewPage from './OverviewPage';
import { BrowserRouter as Router } from 'react-router-dom'; // Assuming you're using react-router

describe('<OverviewPage />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(
      <Router>
        <OverviewPage />
      </Router>
    );
  });

  it('animates the StatCard components', () => {
    cy.mount(
      <Router>
        <OverviewPage />
      </Router>
    );

    // Select the StatCard components using the data-cy attribute
    cy.get('[data-cy=stat-card-time-elapsed]').as('timeElapsedCard');
    cy.get('[data-cy=stat-card-ideal-time]').as('idealTimeCard');
    cy.get('[data-cy=stat-card-tabs-open]').as('tabsOpenCard');
    cy.get('[data-cy=stat-card-data-used]').as('dataUsedCard');

    // Verify initial state (before animation)
    cy.get('@timeElapsedCard').should('have.css', 'opacity', '0');
    cy.get('@idealTimeCard').should('have.css', 'opacity', '0');
    cy.get('@tabsOpenCard').should('have.css', 'opacity', '0');
    cy.get('@dataUsedCard').should('have.css', 'opacity', '0');

    // Wait for animation to complete
    cy.wait(1000);

    // Verify final state (after animation)
    cy.get('@timeElapsedCard').should('have.css', 'opacity', '1');
    cy.get('@idealTimeCard').should('have.css', 'opacity', '1');
    cy.get('@tabsOpenCard').should('have.css', 'opacity', '1');
    cy.get('@dataUsedCard').should('have.css', 'opacity', '1');
  });
});