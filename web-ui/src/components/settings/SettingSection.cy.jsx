import React from 'react'
import SettingSection from './SettingSection'

describe('<SettingSection />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<SettingSection />)
  })
})