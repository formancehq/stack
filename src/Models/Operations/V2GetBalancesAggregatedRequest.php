<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use \formance\stack\Utils\SpeakeasyMetadata;
class V2GetBalancesAggregatedRequest
{
    /**
     * $requestBody
     * 
     * @var ?array<string, mixed> $requestBody
     */
	#[SpeakeasyMetadata('request:mediaType=application/json')]
    public ?array $requestBody = null;
    
    /**
     * Name of the ledger.
     * 
     * @var string $ledger
     */
	#[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=ledger')]
    public string $ledger;
    
	#[SpeakeasyMetadata('queryParam:style=form,explode=true,name=pit,dateTimeFormat=Y-m-d\TH:i:s.up')]
    public ?\DateTime $pit = null;
    
    /**
     * Use insertion date instead of effective date
     * 
     * @var ?bool $useInsertionDate
     */
	#[SpeakeasyMetadata('queryParam:style=form,explode=true,name=use_insertion_date')]
    public ?bool $useInsertionDate = null;
    
	public function __construct()
	{
		$this->requestBody = null;
		$this->ledger = "";
		$this->pit = null;
		$this->useInsertionDate = null;
	}
}
