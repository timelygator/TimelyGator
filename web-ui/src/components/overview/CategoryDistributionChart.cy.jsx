import React from 'react';
import CategoryDistributionChart from './CategoryDistributionChart';

describe('<CategoryDistributionChart />', () => {
  it('renders and animates correctly', () => {
    // Mount the component
    cy.mount(<CategoryDistributionChart />);

    // Check if the motion.div is initially hidden
    cy.get('div.bg-gray-800')
      .should('have.css', 'opacity', '0')
      .and('have.css', 'transform', 'matrix(1, 0, 0, 1, 0, 20)');

    // Wait for the animation to complete
    cy.wait(500); // Adjust the wait time if necessary

    // Check if the motion.div is visible after animation
    cy.get('div.bg-gray-800')
      .should('have.css', 'opacity', '1')
      .and('have.css', 'transform', 'none');
  });
});