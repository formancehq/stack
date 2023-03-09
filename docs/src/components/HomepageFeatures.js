import React from 'react';
import clsx from 'clsx';
import styles from './HomepageFeatures.module.css';
import Icon from "@material-ui/core/Icon";

const FeatureList = [
  {
    title: 'Documentation',
    icon: 'menu_book',
    description: (
      <>
          You will find all the documentation to use Formance Cloud or the OpenSource Ledger.
      </>
    ),
  },
  {
    title: 'Guides',
    icon: 'smart_display',
    description: (
      <>
          You will find concrete examples of how Formance Cloud and the OpenSource Ledger are used
      </>
    ),
  },
  {
    title: 'API',
      icon: 'settings',
      description: (
      <>
          You will find all the documentation for the Formance Cloud APIs as well as for the OpenSource Ledger
      </>
    ),
  },
];

function Feature({icon, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
          <Icon alt={title}>{icon}</Icon>
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
