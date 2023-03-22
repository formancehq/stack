<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;


class GetWalletSummaryResponse
{
	
    public string $contentType;
    
    /**
     * Wallet summary
     * 
     * @var ?\formance\stack\Models\Shared\GetWalletSummaryResponse $getWalletSummaryResponse
     */
	
    public ?\formance\stack\Models\Shared\GetWalletSummaryResponse $getWalletSummaryResponse = null;
    
	
    public int $statusCode;
    
	
    public ?\Psr\Http\Message\ResponseInterface $rawResponse = null;
    
    /**
     * Error
     * 
     * @var ?\formance\stack\Models\Shared\WalletsErrorResponse $walletsErrorResponse
     */
	
    public ?\formance\stack\Models\Shared\WalletsErrorResponse $walletsErrorResponse = null;
    
	public function __construct()
	{
		$this->contentType = "";
		$this->getWalletSummaryResponse = null;
		$this->statusCode = 0;
		$this->rawResponse = null;
		$this->walletsErrorResponse = null;
	}
}
