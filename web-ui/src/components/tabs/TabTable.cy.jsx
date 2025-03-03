import React from 'react';
import TabTable from './TabTable';
import { BrowserRouter as Router } from 'react-router-dom'; // Assuming you're using react-router

describe('<TabTable />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(
      <Router>
        <TabTable />
      </Router>
    );
  });

  it('filters the tabs based on search input', () => {
    cy.mount(
      <Router>
        <TabTable />
      </Router>
    );

    // Select the search input using the data-cy attribute
    cy.get('[data-cy=search-input]').as('searchInput');

    // Verify initial state (all tabs are visible)
    cy.get('tbody tr').should('have.length', 5);

    // Search for "YouTube"
    cy.get('@searchInput').type('YouTube');
    cy.get('tbody tr').should('have.length', 3);

    // Search for "Git"
    cy.get('@searchInput').clear().type('Git');
    cy.get('tbody tr').should('have.length', 2);

    // Search for "Instagram"
    cy.get('@searchInput').clear().type('Instagram');
    cy.get('tbody tr').should('have.length', 1);

    // Search for "Non-existent"
    cy.get('@searchInput').clear().type('Non-existent');
    cy.get('tbody tr').should('have.length', 0);
  });
});