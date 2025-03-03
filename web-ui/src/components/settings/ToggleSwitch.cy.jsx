import React from 'react';
import ToggleSwitch from './ToggleSwitch';

describe('<ToggleSwitch />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<ToggleSwitch label="Test Toggle" isOn={false} onToggle={() => {}} />);
  });

  it('toggles the switch', () => {
    const onToggleSpy = cy.spy().as('onToggleSpy');
    cy.mount(<ToggleSwitch label="Test Toggle" isOn={false} onToggle={onToggleSpy} />);

    // Select the toggle switch using the data-cy attribute
    cy.get('[data-cy=toggle-switch]').as('toggleSwitch');

    // Verify initial state
    cy.get('@toggleSwitch').should('have.class', 'bg-gray-600');

    // Click to toggle on
    cy.get('@toggleSwitch').click();
    cy.get('@toggleSwitch').should('have.class', 'bg-indigo-600');
    cy.get('@onToggleSpy').should('have.been.calledOnce');

    // Click to toggle off
    cy.get('@toggleSwitch').click();
    cy.get('@toggleSwitch').should('have.class', 'bg-gray-600');
    cy.get('@onToggleSpy').should('have.been.calledTwice');
  });
});