import React from 'react'
import SettingsPage from './SettingsPage'

describe('<SettingsPage />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<SettingsPage />)
  })
})