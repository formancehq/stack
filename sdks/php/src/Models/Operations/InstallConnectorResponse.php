<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;


class InstallConnectorResponse
{
    /**
     * OK
     * 
     * @var ?\formance\stack\Models\Shared\ConnectorResponse $connectorResponse
     */
	
    public ?\formance\stack\Models\Shared\ConnectorResponse $connectorResponse = null;
    
	
    public string $contentType;
    
	
    public int $statusCode;
    
	
    public ?\Psr\Http\Message\ResponseInterface $rawResponse = null;
    
	public function __construct()
	{
		$this->connectorResponse = null;
		$this->contentType = "";
		$this->statusCode = 0;
		$this->rawResponse = null;
	}
}
