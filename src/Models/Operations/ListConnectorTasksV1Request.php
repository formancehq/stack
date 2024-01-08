<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use \formance\stack\Utils\SpeakeasyMetadata;
class ListConnectorTasksV1Request
{
    /**
     * The name of the connector.
     * 
     * @var \formance\stack\Models\Shared\Connector $connector
     */
	#[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=connector')]
    public \formance\stack\Models\Shared\Connector $connector;
    
    /**
     * The connector ID.
     * 
     * @var string $connectorId
     */
	#[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=connectorId')]
    public string $connectorId;
    
    /**
     * Parameter used in pagination requests. Maximum page size is set to 15.
     * 
     * Set to the value of next for the next page of results.
     * Set to the value of previous for the previous page of results.
     * No other parameters can be set when this parameter is set.
     * 
     * 
     * @var ?string $cursor
     */
	#[SpeakeasyMetadata('queryParam:style=form,explode=true,name=cursor')]
    public ?string $cursor = null;
    
    /**
     * The maximum number of results to return per page.
     * 
     * 
     * 
     * @var ?int $pageSize
     */
	#[SpeakeasyMetadata('queryParam:style=form,explode=true,name=pageSize')]
    public ?int $pageSize = null;
    
	public function __construct()
	{
		$this->connector = \formance\stack\Models\Shared\Connector::Stripe;
		$this->connectorId = "";
		$this->cursor = null;
		$this->pageSize = null;
	}
}
