<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Models\Shared;
use formance\stack\Utils\SpeakeasyMetadata;
class UninstallConnectorV1Request
{
    /**
     * The name of the connector.
     *
     * @var Shared\Connector $connector
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=connector')]
    public Shared\Connector $connector;

    /**
     * The connector ID.
     *
     * @var string $connectorId
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=connectorId')]
    public string $connectorId;

    /**
     * @param  ?Shared\Connector  $connector
     * @param  ?string  $connectorId
     */
    public function __construct(?Shared\Connector $connector = null, ?string $connectorId = null)
    {
        $this->connector = $connector;
        $this->connectorId = $connectorId;
    }
}